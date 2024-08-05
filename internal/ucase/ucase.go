package ucase

import (
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/healthcheck"
)

type UseCase struct {
	HealthCheck *healthcheck.Ucase
}

func NewUseCase(dep *contract.UseCaseDeps) *UseCase {
	return &UseCase{
		HealthCheck: healthcheck.NewUcase(dep),
	}
}
