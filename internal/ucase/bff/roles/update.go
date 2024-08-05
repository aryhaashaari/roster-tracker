// Package roles
// Automatic generated
package roles

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/presentations"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/repositories"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/validator"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"

	ucase "gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
)

type roleUpdate struct {
	repo repositories.Roler
}

// NewRole new instance
func NewRoleUpdate(repo repositories.Roler) ucase.UseCase {
	return &roleUpdate{repo: repo}
}

// Serve store role data
func (u *roleUpdate) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.RoleParam
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.update")
		lf    = logger.NewFields(
			logger.EventName("Role"),
		)

		code = mux.Vars(dctx.Request)["code"]
	)

	defer tracer.SpanFinish(ctx)

	param.Code = code

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	vRules := []*validation.FieldRules{

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

	df, err := u.repo.FindOne(ctx, code)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error find data to database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	if df == nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("find data %s in database not found", code), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespDataNotFound)
	}

	af, err := u.repo.Update(ctx, param.Email, param.Code)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error update data to database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success update data to database with affected rows %d", af), lf...)
	return *appctx.NewResponse().
		WithMsgKey(consts.RespSuccess)
}
