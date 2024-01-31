package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/pprof"
	"users/config"
	"users/internal/api"
)

func NewRestApi(cfg *config.Config, logger *zerolog.Logger, userService api.UserService) error {
	// инициализируем хэндлер
	userRestHandler := NewUserHandler(logger, userService)

	router := mux.NewRouter()

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

	// зарегистрировать нового пользователя
	router.HandleFunc("/users/register", userRestHandler.RegisterUser).Methods(http.MethodPost)
	// получить пользователя по айди
	router.HandleFunc("/users/{id:[0-9]+}", userRestHandler.GetGetUserById).Methods(http.MethodGet)
	// обновить пользователя
	router.HandleFunc("/users/update", userRestHandler.UpdateUser).Methods(http.MethodPut)
	// обновить пароль
	router.HandleFunc("/users/update-password", userRestHandler.UpdatePassword).Methods(http.MethodPut)
	// удалить пользователя
	router.HandleFunc("/users/delete/{id:[0-9]+}", userRestHandler.DeleteUser).Methods(http.MethodDelete)

	appAddr := fmt.Sprintf("%s:%s", cfg.App.AppHost, cfg.App.AppPort)
	logger.Info().Msgf("running REST server at '%s'", appAddr)
	if err := http.ListenAndServe(appAddr, router); err != nil {
		return fmt.Errorf("[NewRestApi] listen and serve: %w", err)
	}

	return nil
}
