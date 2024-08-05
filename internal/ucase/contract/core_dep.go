package contract

import (
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/business"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/provider"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/repositories"
)

// Core contains all core dependencies, never edit this
type Core struct {
	Cfg      *appctx.Config
	Business *business.Business
	Repo     *repositories.Repo
	//WorkerCl workerx.WorkerClient
	Provider *provider.Provider
}
