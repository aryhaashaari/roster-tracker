// Package middleware
package middleware

type Middleware struct {
	ExampleMiddleware MiddlewareFuncV2
	SpanLog           MiddlewareFuncV2
}

func NewMiddleware(deps *Deps) *Middleware {
	return &Middleware{
		ExampleMiddleware: Example(),
		SpanLog:           SpanLog(),
	}
}
