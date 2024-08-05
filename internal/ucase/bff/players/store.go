// Package players
// Automatic generated
package players

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

type player struct {
	repo repositories.Playerer
}

// NewPlayer new instance
func NewPlayer(repo repositories.Playerer) ucase.UseCase {
	return &player{repo: repo}
}

// Serve store player data
func (u *player) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.PlayerParam
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.create")
		lf    = logger.NewFields(
			logger.EventName("Player"),
		)
	)

	defer tracer.SpanFinish(ctx)

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	vRules := []*validation.FieldRules{

		validation.Field(&param.PlayerName,
			validation.Required,
			validation.Length(1, 0),
			validator.ValidRegex(validator.AlphaDashSpace),
		),

		validation.Field(&param.Position,
			validation.Required,
			validation.Length(1, 0),
			validator.ValidRegex(validator.Alpha),
		),

		// validation.Field(&param.Physique,
		// 	validation.Required,
		// 	validation.Length(1, 0),
		// 	validator.ValidRegex(validator.AlphaNumeric),
		// ),

		// validation.Field(&param.Stats,
		// 	validation.Required,
		// 	validation.Length(1, 0),
		// 	validator.ValidRegex(validator.AlphaNumeric),
		// ),
	}

	err = validation.ValidateStruct(&param, vRules...)

	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("validation error %v", err), lf...)
		return *appctx.NewResponse().
			WithMsgKey(consts.RespValidationError).
			WithError(validator.ExtractMessageToSliceMap(err))
	}

	af, err := u.repo.Store(ctx, param.PlayerName, param.Position, param.Physique, param.StatsList)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.WarnWithContext(ctx, fmt.Sprintf("store data to database error: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	lf.Append(logger.Any("affected_rows", af))

	logger.InfoWithContext(ctx, fmt.Sprintf("success store data to database"), lf...)
	return *appctx.NewResponse().
		WithMsgKey(consts.RespSuccess)
}
