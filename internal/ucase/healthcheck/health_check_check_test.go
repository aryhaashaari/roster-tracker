package healthcheck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/healthcheck"
)

func TestHealthCheck_Serve(t *testing.T) {
	svc := healthcheck.NewCheck(&contract.UseCaseDeps{}, "")

	t.Run("test health check", func(t *testing.T) {
		result := svc.Serve(&appctx.Data{})

		assert.Equal(t, appctx.Response{
			Code:    consts.CodeSuccess,
			Message: "ok",
		}, result)
	})
}
