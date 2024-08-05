// Package router
package router

import (
	"net/http"
	"runtime/debug"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/bootstrap"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/business"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/handler"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/middleware"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/provider"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/repositories"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/msgx"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/routerkit"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/bff/players"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/bff/roles"
	ucaseContract "gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/jwt"
)

type router struct {
	config           *appctx.Config
	router           *routerkit.Router
	globalMiddleware []middleware.MiddlewareFuncV2
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)
	bootstrap.RegistrySnowflake()

	return &router{
		config:           cfg,
		router:           routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
		globalMiddleware: make([]middleware.MiddlewareFuncV2, 0),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFuncV2) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get(consts.HeaderLanguageKey)
		if !msgx.HaveLang(consts.RespOK, lang) {
			lang = rtr.config.App.DefaultLang
			r.Header.Set(consts.HeaderLanguageKey, lang)
		}

		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Code: consts.CodeInternalServerError,
				}

				res.WithLang(lang)
				logger.Error(logger.MessageFormat("got panic error: %v", err))
				logger.Error(logger.MessageFormat("got panic: stack trace: %v", string(debug.Stack())))

				rtr.response(w, res)

				return
			}
		}()

		// validate middleware
		adapter := func(w http.ResponseWriter, r *http.Request) appctx.Response {
			resp := hfn(r, svc, rtr.config)
			resp.WithLang(lang)
			return resp
		}

		global := rtr.globalMiddleware

		middlewares := append(global, mdws...)

		fn := middleware.Wrap(adapter, middlewares...)

		res := fn(w, r)

		rtr.response(w, res)
	}
}

func (rtr *router) UseGlobalMiddlewares(mdws ...middleware.MiddlewareFuncV2) {
	rtr.globalMiddleware = mdws
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {
	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
	resp.Generate()
	w.WriteHeader(resp.Code)
	w.Write(resp.Byte())
	return
}

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	rtr.router.NotFoundHandler = http.HandlerFunc(middleware.NotFound)

	//root := rtr.router.PathPrefix("/").Subrouter()
	//root.Use(middleware.ProcessIdInjector)
	//in := root.PathPrefix("/in/").Subrouter()
	//liveness := root.PathPrefix("/").Subrouter()
	//inV1 := in.PathPrefix("/v1/").Subrouter()

	//_ = inV1

	// all bootstrap init
	bootstrap.RegistryOpenTracing(rtr.config)

	// Create databaase
	db := bootstrap.RegistryDatabase(rtr.config.WriteDB)

	// all deps init
	roleRepo := repositories.NewRole(db)
	playerRepo := repositories.NewPlayer(db)
	repo := repositories.NewRepo(&repositories.Deps{})
	providerCl := provider.NewProvider(&provider.Deps{})
	businessLogic := business.NewBusiness(&business.Deps{})

	uc := ucase.NewUseCase(&ucaseContract.UseCaseDeps{
		Core: ucaseContract.Core{
			Cfg:      rtr.config,
			Business: businessLogic,
			Repo:     repo,
			Provider: providerCl,
		},
	})

	uci := jwt.NewUcase(&ucaseContract.UseCaseDeps{
		Core: ucaseContract.Core{
			Cfg:      rtr.config,
			Business: businessLogic,
			Repo:     repo,
			Provider: providerCl,
		},
	})

	m := middleware.NewMiddleware(&middleware.Deps{})

	// global middleware
	rtr.UseGlobalMiddlewares(
		m.SpanLog,
	)

	// Use Case
	createRole := roles.NewRole(roleRepo)
	updateRole := roles.NewRoleUpdate(roleRepo)
	createPlayer := players.NewPlayer(playerRepo)
	getPlayers := players.NewPlayerList(playerRepo)
	getSinglePlayer := players.NewSinglePlayer(playerRepo)
	updatePlayer := players.NewPlayerUpdate(playerRepo)
	deletePlayer := players.NewPlayerDelete(playerRepo)

	// router
	rtr.router.HandleFunc("/liveness", rtr.handle(
		handler.HttpRequest,
		uc.HealthCheck.Check,
		m.ExampleMiddleware,
	)).Methods(http.MethodGet)

	rtr.router.HandleFunc("/login", rtr.handle(
		handler.HttpRequest,
		uci.Check,
		m.ExampleMiddleware,
	)).Methods(http.MethodGet)

	// router
	rtr.router.HandleFunc("/role/register", rtr.handle(
		handler.HttpRequest,
		createRole,
	)).Methods(http.MethodPost)

	// router
	rtr.router.HandleFunc("/role/{code}", rtr.handle(
		handler.HttpRequest,
		updateRole,
	)).Methods(http.MethodPatch)

	// router
	rtr.router.HandleFunc("/addPlayer", rtr.handle(
		handler.HttpRequest,
		createPlayer,
	)).Methods(http.MethodPost)

	// router
	rtr.router.HandleFunc("/getPlayers", rtr.handle(
		handler.HttpRequest,
		getPlayers,
	)).Methods(http.MethodGet)

	// router
	rtr.router.HandleFunc("/getPlayer/{id}", rtr.handle(
		handler.HttpRequest,
		getSinglePlayer,
	)).Methods(http.MethodGet)

	// router
	rtr.router.HandleFunc("/updatePlayer/{id}", rtr.handle(
		handler.HttpRequest,
		updatePlayer,
	)).Methods(http.MethodPut)

	rtr.router.HandleFunc("/deletePlayer/{id}", rtr.handle(
		handler.HttpRequest,
		deletePlayer,
	)).Methods(http.MethodDelete)

	// this route for example rest, please delete
	// example list
	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    el,
	//)).Methods(http.MethodGet)

	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    ec,
	//)).Methods(http.MethodPost)

	//inV1.HandleFunc("/example/{id:[0-9]+}", rtr.handle(
	//    handler.HttpRequest,
	//    ed,
	//)).Methods(http.MethodDelete)

	// mount router here

	return rtr.router

}
