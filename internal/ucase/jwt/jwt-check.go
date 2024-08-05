package jwt

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/middleware"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
)

type Check struct {
	dep    *contract.UseCaseDeps
	entity string
}

func NewCheck(dep *contract.UseCaseDeps, entity string) contract.UseCase {
	return &Check{
		dep:    dep,
		entity: entity,
	}
}

func (c *Check) Serve(data *appctx.Data) appctx.Response {

	responder := appctx.NewResponse().WithEntity(c.entity)

	var err error = nil
	lf := logger.NewFields(logger.Field{
		Key:   "key-1",
		Value: "test key",
	})

	// example using responder translation with extender and logger field
	if err != nil {
		return *responder.ResponderTranslationError(data.Request.Context(), err, responder.WithLoggerField(lf...),
			responder.WithExtender(func(ctx context.Context, err error, fl ...logger.Field) *appctx.Response {

				switch err.(type) {
				case CustomError:
					return responder.WithCode(http.StatusUnauthorized).WithMessage(err)
				}

				return nil
			}))
	}

	valid, err := middleware.ValidateHMAC(data.Request, data.Config.App.SecretKey)

	if err != nil {
		return *responder.
			WithCode(consts.CodeInternalServerError).
			WithMessage("Internal Server Error")
	}

	if !valid {
		return *responder.
			WithCode(consts.CodeAuthenticationFailure).
			WithMessage("Invalid Signature")
	}
	token, err := middleware.CreateJwtToken(data.Config.App.SecretKey)

	if err != nil {
		return *responder.
			WithCode(consts.CodeInternalServerError).
			WithMessage("Internal Server Error")
	}

	response := map[string]interface{}{
		"expired_at": time.Now().Add(time.Hour * 2),
		"token":      token,
		"type":       "Bearer",
	}

	meta := map[string]interface{}{
		"transaction_id": uuid.New(),
	}

	return *responder.
		WithCode(consts.CodeSuccess).
		WithStatus(consts.StatusSuccess).
		WithState("generateTokenSuccess").
		WithMessage("Generate Token Success").
		WithData(meta).
		WithData(response)

}

type CustomError string

func (c CustomError) Error() string {
	return string(c)
}
