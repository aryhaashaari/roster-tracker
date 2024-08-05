// Package handler
package handler

import (
	"context"
	"net/http"
	"time"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
)

// HttpRequest handler func wrapper
func HttpRequest(request *http.Request, svc contract.UseCase, conf *appctx.Config) appctx.Response {
	ctx := context.WithValue(request.Context(), consts.CtxLang, request.Header.Get(consts.HeaderLanguageKey))

	req := request.WithContext(ctx)

	data := &appctx.Data{
		Request:     req,
		Config:      conf,
		ServiceType: consts.ServiceTypeHTTP,
		RequestID:   req.Header.Get(`X-Request-Id`),
	}

	t, ok := ctx.Value(consts.CtxtStartTime).(time.Time)
	if ok {
		data.StartTime = t
	}

	return svc.Serve(data)
}
