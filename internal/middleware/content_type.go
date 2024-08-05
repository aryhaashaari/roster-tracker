// Package middleware
package middleware

import (
	"net/http"
	"strings"
	"fmt"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
)

// ValidateContentType header
func ValidateContentType(r *http.Request, conf *appctx.Config) int {

	if ct := strings.ToLower(r.Header.Get(`Content-Type`)) ; ct != `application/json` {
		logger.Warn(fmt.Sprintf("[middleware] invalid content-type %s", ct ))

		return consts.CodeBadRequest
	}


	return consts.CodeSuccess
}