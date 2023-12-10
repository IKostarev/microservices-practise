package app

import (
	"fmt"
	"gateway/config"
	"gateway/internal/handlers"
	"gateway/internal/service"
	"gateway/pkg/clients/todos_client"
	"gateway/pkg/clients/users_client"
	"gateway/pkg/logging"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

type App struct {
	cfg            *config.Config
	logger         *zerolog.Logger
	router         *mux.Router
	gatewayService handlers.GatewayService
}

func NewApp(cfg *config.Config) (*App, error) {

	todosClient := todos_client.NewTodosClient()
	usersClient := users_client.NewUsersClient()

	gatewayService := service.NewGatewayService(&cfg.JWT, todosClient, usersClient)

	logger := logging.NewLogger(cfg.Logging)

	return &App{
		cfg:            cfg,
		logger:         logger,
		gatewayService: gatewayService,
	}, nil
}

func (a *App) RunAPI() {
	gatewayHandler := handlers.NewGatewayHandler(a.logger, a.gatewayService)

	a.router = mux.NewRouter()
	a.router.Use(
		handlers.ValidateTokenMiddleware(
			&a.cfg.JWT,
			[]string{
				"/api/v1/users/login",
				"/api/v1/users/register",
			},
		),
	)

	usersV1Router := a.router.PathPrefix("/api/v1/users").Subrouter()
	usersV1Router.HandleFunc("/register", gatewayHandler.RegisterUser).Methods(http.MethodPost)
	usersV1Router.HandleFunc("/{id:[0-9]+}", gatewayHandler.GetGetUserById).Methods(http.MethodGet)
	usersV1Router.HandleFunc("/update", gatewayHandler.UpdateUser).Methods(http.MethodPut)
	usersV1Router.HandleFunc("/update-password", gatewayHandler.UpdatePassword).Methods(http.MethodPut)
	usersV1Router.HandleFunc("/delete/{id:[0-9]+}", gatewayHandler.DeleteUser).Methods(http.MethodDelete)
	usersV1Router.HandleFunc("/login", gatewayHandler.UserLogin).Methods(http.MethodPost)
	usersV1Router.HandleFunc("/refresh", gatewayHandler.Refresh).Methods(http.MethodPost)

	todosV1Router := a.router.PathPrefix("/api/v1/todos").Subrouter()
	todosV1Router.HandleFunc("/", gatewayHandler.CreateToDoHandler).Methods(http.MethodPost)
	todosV1Router.HandleFunc("/batch", gatewayHandler.GetToDosHandler).Methods(http.MethodGet)
	todosV1Router.HandleFunc("/{id}", gatewayHandler.GetToDoHandler).Methods(http.MethodGet)
	todosV1Router.HandleFunc("/{id}", gatewayHandler.UpdateToDoHandler).Methods(http.MethodPut)
	todosV1Router.HandleFunc("/{id}", gatewayHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	// запустить вебсервер по адресу, передать в него роутер
	appAddr := fmt.Sprintf("%s:%s", a.cfg.App.AppHost, a.cfg.App.AppPort)
	a.logger.Info().Msgf("running server at '%s'", appAddr)
	http.ListenAndServe(appAddr, a.router)
}
