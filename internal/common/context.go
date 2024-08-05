package common

import (
	"context"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
)

func RequestIdFromContext(ctx context.Context) string {
	v, ok := ctx.Value(consts.ContextRequestID).(string)
	if !ok {
		return "-"
	}

	return v
}
