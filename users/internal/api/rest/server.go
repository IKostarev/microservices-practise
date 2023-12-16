package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"users/config"
	"users/internal/api"
)

func NewRestApi(cfg *config.Config, logger *zerolog.Logger, userService api.UserService) error {
	// инициализируем хэндлер
	userRestHandler := NewUserHandler(logger, userService)

	router := mux.NewRouter()

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
