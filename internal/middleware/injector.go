package middleware

import (
	"context"
	"net/http"
	"time"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/util"
)

func ProcessIdInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		regid := r.Header.Get(`X-Request-Id`)
		if regid == "" {
			regid = util.GenerateReferenceID("PID-")
			r.Header.Set(`X-Request-Id`, regid)
		}
		w.Header().Set(`X-Request-Id`, regid)

		ctx := context.WithValue(r.Context(), consts.ContextRequestID, regid)
        ctx = context.WithValue(ctx, consts.ContextStartTime, time.Now())
        ctx = context.WithValue(ctx, consts.ContextRequestIp, util.IPFromRequest(r))
        ctx = context.WithValue(ctx, consts.ContextRequestPath, r.URL.Path)
        ctx = context.WithValue(ctx, consts.ContextRequestMethod, r.Method)


		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
