package httpclientx

import (
	"bytes"
	"github.com/opentracing/opentracing-go"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"
	"io"
	"net/http"
	"strconv"
)

func Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	tracer.AddSpanTag(ctx,
		tracer.SpanTag{Key: "http.request.headers.*", Value: req.Header},
		tracer.SpanTag{Key: "http.url", Value: req.URL.Path},
		tracer.SpanTag{Key: "http.method", Value: req.Method},
	)

	if req.Body != nil {
		// Re-usable request body for logging
		requestBody, _ := io.ReadAll(req.Body)
		req.Body.Close() // must close
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		tracer.AddSpanTag(ctx,
			tracer.SpanTag{Key: "http.request.body", Value: string(requestBody)},
		)
	}

	//Add X-Request-ID to the request Header
	//state, valid := ctx.Value(consts.CtxRequestState).(common.RequestState)
	//if valid {
	//	req.Header.Set(consts.HeaderXRequestID, state.ID)
	//}

	// Attempt to join a trace by getting trace context from the headers.
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	span := opentracing.SpanFromContext(ctx)

	// Transmit the span's TraceContext as HTTP headers on our
	// outbound request.
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		carrier)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Re-usable response body for logging
	responseBody, _ := io.ReadAll(res.Body)
	res.Body.Close() // must close
	res.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	tracer.AddSpanTag(ctx,
		tracer.SpanTag{Key: "http.status_code", Value: strconv.Itoa(res.StatusCode)},
		tracer.SpanTag{Key: "http.response.headers.*", Value: res.Header},
		tracer.SpanTag{Key: "http.response.body", Value: string(responseBody)},
	)

	return res, nil
}
