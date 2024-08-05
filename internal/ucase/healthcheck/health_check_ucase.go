package healthcheck

import "gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"

type Ucase struct {
	Check contract.UseCase
}

func NewUcase(dep *contract.UseCaseDeps) *Ucase {
	return &Ucase{
		Check: NewCheck(dep, "HealthCheck"),
	}
}
