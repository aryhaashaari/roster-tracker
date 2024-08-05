// Package players
// Automatic generated
package players

import (
	"fmt"

	"github.com/gorilla/mux"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/dto"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/presentations"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/repositories"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"

	ucase "gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
)

type playerSingle struct {
	repo repositories.Playerer
}

func NewSinglePlayer(repo repositories.Playerer) ucase.UseCase {
	return &playerSingle{repo: repo}
}

// Serve Player list data
func (u *playerSingle) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.PlayerQuery
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.single_player")
		lf    = logger.NewFields(
			logger.EventName("singlePlayer"),
		)

		id = mux.Vars(dctx.Request)["id"]
	)
	defer tracer.SpanFinish(ctx)

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	// vRules := []*validation.FieldRules{
	// 	validation.Field(&param.Page, validation.Min(int64(1))),
	// 	validation.Field(&param.Limit, validation.Min(int64(1))),
	// 	validation.Field(&param.StartDate, validation.Required.When(param.EndDate != ""), validator.ValidDateTime()),
	// 	validation.Field(&param.EndDate, validation.Required.When(param.StartDate != ""), validator.ValidDateTime()),
	// }

	// err = validation.ValidateStruct(&param, vRules...)

	// if err != nil {
	// 	logger.WarnWithContext(ctx, fmt.Sprintf("validation error %v", err), lf...)
	// 	return *appctx.NewResponse().
	// 		WithMsgKey(consts.RespValidationError).
	// 		WithError(validator.ExtractMessageToSliceMap(err))
	// }

	// param.Limit = common.LimitDefaultValue(param.Limit)
	// param.Page = common.PageDefaultValue(param.Page)

	dr, err := u.repo.FindOne(ctx, id)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error find data to database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	if dr == nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("find data %s in database not found", id), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespDataNotFound)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success fetch players to database"), lf...)
	return *appctx.NewResponse().
		WithMsgKey(consts.RespSuccess).
		WithData(dto.DetailedPlayerToResponse(dr))

}
