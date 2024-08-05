// Package players
// Automatic generated
package players

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

type playerDelete struct {
	repo repositories.Playerer
}

// NewPlayer new instance
func NewPlayerDelete(repo repositories.Playerer) ucase.UseCase {
	return &playerDelete{repo: repo}
}

// Serve store player data
func (u *playerDelete) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.PlayerParam
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.delete")
		lf    = logger.NewFields(
			logger.EventName("Player"),
		)

		id = mux.Vars(dctx.Request)["id"]
	)

	defer tracer.SpanFinish(ctx)

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	vRules := []*validation.FieldRules{

		// validation.Field(&param.PlayerName,
		// 	validation.Required,
		// 	validation.Length(1, 0),
		// 	validator.ValidRegex(validator.AlphaNumeric),
		// ),

		// validation.Field(&param.Position,
		// 	validation.Required,
		// 	validation.Length(1, 0),
		// 	validator.ValidRegex(validator.AlphaNumeric),
		// ),

		// validation.Field(&param.Physique,
		// 	validation.Required,
		// 	validation.Length(1, 0),
		// 	validator.ValidRegex(validator.AlphaNumeric),
		// ),

		// validation.Field(&param.StatsList,
		// 	validation.Required,
		// 	validation.Length(1, 0),
		// 	validator.ValidRegex(validator.AlphaNumeric),
		// ),

		// validation.Field(&param.AvgStats,
		//     validation.Required,
		//     validation.Length(1, 0),
		//     validator.ValidRegex(validator.AlphaNumeric),
		// ),
	}

	err = validation.ValidateStruct(&param, vRules...)

	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("validation error %v", err), lf...)
		return *appctx.NewResponse().
			WithMsgKey(consts.RespValidationError).
			WithError(validator.ExtractMessageToSliceMap(err))
	}

	df, err := u.repo.FindOne(ctx, id)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error find data to database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	if df == nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("find data %s in database not found", id), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespDataNotFound)
	}

	af, err := u.repo.Delete(ctx, id)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error deleting data from database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success deleted data from database with affected rows %d", af), lf...)
	return *appctx.NewResponse().
		WithMsgKey(consts.RespSuccess)
}
