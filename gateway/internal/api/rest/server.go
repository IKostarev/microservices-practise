package rest

import (
	"fmt"
	"gateway/config"
	rest "gateway/internal/api"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

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
				"/api/v1/users/login",
				"/api/v1/users/register",
			},
		),
	)

	usersV1Router := router.PathPrefix("/api/v1/users").Subrouter()
	usersV1Router.HandleFunc("/register", gatewayHandler.RegisterUser).Methods(http.MethodPost)
	usersV1Router.HandleFunc("/{id:[0-9]+}", gatewayHandler.GetGetUserById).Methods(http.MethodGet)
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
