package rest

import (
	"fmt"
	"gateway/config"
	rest "gateway/internal/api"
	_ "gateway/internal/api/docs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/http/pprof"
)

// @title           ToDo Gateway API
// @version         1.0
// @description     This service is Gateway API for all microservices of ToDo service
// @host            localhost:3000
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func RunREST(
	cfg *config.Config,
	logger *zerolog.Logger,
	gatewayService rest.GatewayService,
) error {
	gatewayHandler := NewGatewayHandler(logger, gatewayService)

	router := mux.NewRouter()
	router.Use(
		ValidateTokenMiddleware(
			&cfg.JWT,
			[]string{
				"/debug/pprof",
				"/docs/swagger",
				"/api/v1/users/login",
				"/api/v1/users/register",
			},
		),
	)

	router.PathPrefix("/docs/swagger").Handler(httpSwagger.WrapHandler)

	debugRouter := router.PathPrefix("/debug/pprof").Subrouter()

	debugRouter.HandleFunc("/", pprof.Index).Methods(http.MethodGet)
	debugRouter.HandleFunc("/cmdline", pprof.Cmdline).Methods(http.MethodGet)
	debugRouter.HandleFunc("/profile", pprof.Profile).Methods(http.MethodGet)
	debugRouter.HandleFunc("/symbol", pprof.Symbol).Methods(http.MethodGet)
	debugRouter.HandleFunc("/trace", pprof.Trace).Methods(http.MethodGet)
	debugRouter.Handle("/goroutine", pprof.Handler("goroutine")).Methods(http.MethodGet)
	debugRouter.Handle("/heap", pprof.Handler("heap")).Methods(http.MethodGet)
	debugRouter.Handle("/threadcreate", pprof.Handler("threadcreate")).Methods(http.MethodGet)
	debugRouter.Handle("/block", pprof.Handler("block")).Methods(http.MethodGet)

	usersV1Router := router.PathPrefix("/api/v1/users").Subrouter()
	usersV1Router.HandleFunc("/register", gatewayHandler.RegisterUser).Methods(http.MethodPost)
	usersV1Router.HandleFunc("/{id:[0-9]+}", gatewayHandler.GetUserById).Methods(http.MethodGet)
	usersV1Router.HandleFunc("/update", gatewayHandler.UpdateUser).Methods(http.MethodPut)
	usersV1Router.HandleFunc("/update-password", gatewayHandler.UpdatePassword).Methods(http.MethodPut)
	usersV1Router.HandleFunc("/delete/{id:[0-9]+}", gatewayHandler.DeleteUser).Methods(http.MethodDelete)
	usersV1Router.HandleFunc("/login", gatewayHandler.UserLogin).Methods(http.MethodPost)
	usersV1Router.HandleFunc("/refresh", gatewayHandler.Refresh).Methods(http.MethodPost)

	todosV1Router := router.PathPrefix("/api/v1/todos").Subrouter()
	todosV1Router.HandleFunc("/", gatewayHandler.CreateToDoHandler).Methods(http.MethodPost)
	todosV1Router.HandleFunc("/batch", gatewayHandler.GetToDosHandler).Methods(http.MethodGet)
	todosV1Router.HandleFunc("/{id}", gatewayHandler.GetToDoHandler).Methods(http.MethodGet)
	todosV1Router.HandleFunc("/{id}", gatewayHandler.UpdateToDoHandler).Methods(http.MethodPut)
	todosV1Router.HandleFunc("/{id}", gatewayHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	// запустить вебсервер по адресу, передать в него роутер
	appAddr := fmt.Sprintf("%s:%s", cfg.App.AppHost, cfg.App.AppPort)
	logger.Info().Msgf("running server at '%s'", appAddr)
	err := http.ListenAndServe(appAddr, router)
	if err != nil {
		return fmt.Errorf("[RunREST] listen and serve: %w", err)
	}

	return nil
}
