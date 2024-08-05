// Package roles
// Automatic generated
package roles

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/presentations"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/repositories"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/validator"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"

	ucase "gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
)

type role struct {
	repo repositories.Roler
}

// NewRole new instance
func NewRole(repo repositories.Roler) ucase.UseCase {
	return &role{repo: repo}
}

// Serve store role data
func (u *role) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.RoleParam
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.create")
		lf    = logger.NewFields(
			logger.EventName("Role"),
		)
	)

	defer tracer.SpanFinish(ctx)

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	vRules := []*validation.FieldRules{

		validation.Field(&param.Email,
			validation.Required,
			validation.Length(1, 50),
			validator.ValidRegex(validator.Email),
		),

		validation.Field(&param.Password,
			validation.Required,
			validation.Length(1, 50),
			validator.ValidRegex(validator.AlphaNumericDash),
		),

		validation.Field(&param.Code,
			validation.Required,
			validation.Length(1, 50),
			validator.ValidRegex(validator.Alpha),
		),
	}

	err = validation.ValidateStruct(&param, vRules...)

	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("validation error %v", err), lf...)
		return *appctx.NewResponse().
			WithMsgKey(consts.RespValidationError).
			WithError(validator.ExtractMessageToSliceMap(err))
	}

	err = u.repo.Store(ctx, param.Email, param.Password, param.Code)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.WarnWithContext(ctx, fmt.Sprintf("store data to database error: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success store data to database"), lf...)
	return *appctx.NewResponse().
		WithMsgKey(consts.RespSuccess)
}
