package middleware

import (
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"net/http"
)

// MiddlewareFunc is contract for middleware and must implement this type for http if need middleware http request
type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, conf *appctx.Config) bool

// FilterFunc is a iterator resolver in each middleware registered
func FilterFunc(w http.ResponseWriter, r *http.Request, conf *appctx.Config, mfs []MiddlewareFunc) bool {
	for _, mf := range mfs {
		return mf(w, r, conf)
	}

	return true
}

type HandleFunc func(w http.ResponseWriter, r *http.Request) appctx.Response

type MiddlewareFuncV2 func(handlerFunc HandleFunc) HandleFunc

func Wrap(handlerFunc HandleFunc, mdws ...MiddlewareFuncV2) HandleFunc {
	fn := handlerFunc

	for i := len(mdws) - 1; i >= 0; i-- {
		fn = mdws[i](fn)
	}

	return fn
}
