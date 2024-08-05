package middleware

import (
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"net/http"
)

func Example() MiddlewareFuncV2 {
	return func(next HandleFunc) HandleFunc {
		return func(writer http.ResponseWriter, r *http.Request) appctx.Response {

			// any pre proccess middleware

			res := next(writer, r)

			//	res contains appctx.Response from next process
			// post middleware

			return res
		}
	}
}
