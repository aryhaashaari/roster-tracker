package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"
	"io"
	"net/http"
	"strings"
)

func InternalHmacMiddleware(conf *appctx.Config) MiddlewareFuncV2 {
	return func(next HandleFunc) HandleFunc {
		return func(writer http.ResponseWriter, r *http.Request) appctx.Response {

			var (
				clientID  = strings.TrimSpace(r.Header.Get("Client-Id"))
				timestamp = strings.TrimSpace(r.Header.Get("Timestamp"))
				signature = strings.TrimSpace(r.Header.Get("Signature"))
				deviceID  = strings.TrimSpace(r.Header.Get("Device-Id"))

				eventName = "Middleware.HMAC"
				ctx       = tracer.SpanStart(r.Context(), eventName)
				lf        = logger.NewFields(
					logger.EventName(eventName),
					logger.String("client_id", clientID),
					logger.String("device_id", deviceID),
				)
				errorReturn = appctx.NewResponse().
						WithCode(http.StatusUnauthorized).
						WithStatus(consts.RespError).
						WithEntity("validateHMAC").
						WithState("validateHMACFailed").
						WithMessage("Validate HMAC Failed")
			)

			defer tracer.SpanFinish(ctx)

			if len(clientID) == 0 || len(deviceID) == 0 || len(timestamp) == 0 || len(signature) == 0 {
				var missingHeaders []string
				if len(clientID) == 0 {
					missingHeaders = append(missingHeaders, "Client-Id")
				}
				if len(deviceID) == 0 {
					missingHeaders = append(missingHeaders, "Device-Id")
				}
				if len(timestamp) == 0 {
					missingHeaders = append(missingHeaders, "Timestamp")
				}
				if len(signature) == 0 {
					missingHeaders = append(missingHeaders, "Signature")
				}
				tracer.SpanError(ctx, errors.New("missing_header"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error missing header %v", eventName, missingHeaders), lf...)
				return appctx.Response{
					Code:    http.StatusBadRequest,
					Message: "invalid request",
					Errors: map[string]interface{}{
						"missing_header": missingHeaders,
					},
				}
			}

			rawRequestBody, errDecode := io.ReadAll(r.Body)
			if errDecode != nil {
				tracer.SpanError(ctx, errors.New("failed to read body"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error read body %v", eventName, errDecode), lf...)
				return *errorReturn.
					WithCode(http.StatusPreconditionFailed).
					WithMessage("failed to read body")
			}
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(rawRequestBody))

			decodedSignature := make([]byte, base64.StdEncoding.DecodedLen(len(signature)))
			xSignature, errDecodeBase64 := base64.StdEncoding.Decode(decodedSignature, []byte(signature))
			if errDecodeBase64 != nil {
				tracer.SpanError(ctx, errors.New("decode signature not found"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error read signature %v", eventName, errDecodeBase64), lf...)
				return appctx.Response{
					Code:    http.StatusBadRequest,
					Message: "invalid request",
					Errors:  errDecodeBase64,
				}
			}

			decodedSignature = decodedSignature[:xSignature]
			if !strings.Contains(string(decodedSignature), ":") && !strings.Contains(string(decodedSignature), "#") {
				tracer.SpanError(ctx, errors.New("decode signature not found"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error signature not valid %v", eventName, errors.New("decode signature not found")), lf...)
				return appctx.Response{
					Code:    http.StatusBadRequest,
					Message: "invalid request",
					Errors:  errors.New("decode signature not found"),
				}
			}

			splitSignature := strings.Split(string(decodedSignature), ":")
			if len(splitSignature) != 2 {
				tracer.SpanError(ctx, errors.New("split signature error"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error split signature %v", eventName, errors.New("split signature error")), lf...)
				return *errorReturn.WithMessage("split signature error")
			}

			trimmedClientID := strings.TrimPrefix(splitSignature[0], "#")
			trimmedSignature := strings.TrimPrefix(splitSignature[1], "#")

			if clientID != trimmedClientID || clientID != conf.App.ClientID {
				tracer.SpanError(ctx, errors.New("client id not valid"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error client-id not valid %v", eventName, errors.New("client id not valid")), lf...)
				return *errorReturn.WithMessage("client id not valid")
			}

			requestBody, errRequestBody := RequestBodyString(rawRequestBody)
			if errRequestBody != nil {
				tracer.SpanError(ctx, errors.New("client not found"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] parse to string no valid %v", eventName, errRequestBody), lf...)
				return *errorReturn.WithMessage("client id not found")
			}

			bodyMD5 := md5.Sum([]byte(requestBody))
			bodyMD5Base64 := make([]byte, base64.StdEncoding.EncodedLen(len(bodyMD5)))
			base64.StdEncoding.Encode(bodyMD5Base64, bodyMD5[:])

			hash := hmac.New(sha256.New, []byte(conf.App.ClientSecret))
			hash.Write([]byte(fmt.Sprintf("%s:%s:%s:%s", timestamp, clientID, r.Method, string(bodyMD5Base64))))
			hmacBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(hash.Sum(nil))))
			base64.StdEncoding.Encode(hmacBase64, hash.Sum(nil))

			if string(hmacBase64) != trimmedSignature {
				tracer.SpanError(ctx, errors.New("hmac not valid"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error hmac signature not valid %v", eventName, errors.New("hmac not valid")), lf...)
				return *errorReturn.WithMessage("hmac not valid")
			}

			generatedSignature := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("#%s:#%s", clientID, string(hmacBase64))))

			if generatedSignature != signature {
				tracer.SpanError(ctx, errors.New("signature not valid"))
				logger.ErrorWithContext(ctx, fmt.Sprintf("[%s] got error signature not valid %v", eventName, errors.New("signature not valid")), lf...)
				return *errorReturn.WithMessage("signature not valid")
			}

			res := next(writer, r)

			return res
		}
	}
}

func ExternalHmacMiddleware(conf *appctx.Config) MiddlewareFuncV2 {
	return func(next HandleFunc) HandleFunc {
		return func(writer http.ResponseWriter, r *http.Request) appctx.Response {

			// connect with provider oauth

			res := next(writer, r)

			return res
		}
	}
}

func RequestBodyString(rawRequestBody []byte) (string, error) {
	switch len(rawRequestBody) {
	case 0:
		emptyString := ""
		return emptyString, nil
	default:
		var requestBodyBuff bytes.Buffer
		errCompactJSON := json.Compact(&requestBodyBuff, rawRequestBody)
		if errCompactJSON != nil {
			emptyString := ""
			return emptyString, errCompactJSON
		}
		result := strings.ReplaceAll(requestBodyBuff.String(), " ", "")
		return result, nil
	}
}
