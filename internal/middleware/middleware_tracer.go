package middleware

import (
	"bytes"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"
	"io"
	"net/http"
)

type customWriter struct {
	http.ResponseWriter
	Body []byte
}

func (c *customWriter) Write(b []byte) (int, error) {
	c.Body = b

	return c.ResponseWriter.Write(b)
}

func SpanLog() MiddlewareFuncV2 {
	return func(handlerFunc HandleFunc) HandleFunc {
		return func(writer http.ResponseWriter, request *http.Request) appctx.Response {
			ctx := tracer.SpanStart(request.Context(), "middleware.trace_request_response")
			defer tracer.SpanFinish(ctx)

			// Re-usable response body for logging
			requestBody, _ := io.ReadAll(request.Body)
			request.Body.Close() // must close
			request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			tracer.AddSpanTag(ctx,
				tracer.NewSpanTag("http.request.headers.*", request.Header),
				tracer.NewSpanTag("http.request.body", string(requestBody)),
				tracer.NewSpanTag("http.request.query_params", request.URL.Query()),
			)

			cw := customWriter{ResponseWriter: writer}

			errHandle := handlerFunc(&cw, request)

			tracer.AddSpanTag(ctx,
				tracer.NewSpanTag("http.response.headers.*", writer.Header()),
				tracer.NewSpanTag("http.response.body", string(cw.Body)),
			)

			return errHandle
		}
	}
}
