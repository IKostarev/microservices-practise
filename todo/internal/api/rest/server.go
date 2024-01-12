package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"todo/config"
	"todo/internal/api"
)

func NewRestApi(cfg *config.Config, logger *zerolog.Logger, todoService api.TodoService) error {
	todoRestHandler := NewTodoHandler(logger, todoService)

	router := mux.NewRouter()

	router.HandleFunc("/todo/", todoRestHandler.CreateToDoHandler).Methods(http.MethodPost)

	router.HandleFunc("/todo/batch", todoRestHandler.GetToDosHandler).Methods(http.MethodGet)

	router.HandleFunc("/todo/{id}", todoRestHandler.GetToDoHandler).Methods(http.MethodGet)

	router.HandleFunc("/todo/{id}", todoRestHandler.UpdateToDoHandler).Methods(http.MethodPut)

	router.HandleFunc("/todo/{id}", todoRestHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	appAddr := fmt.Sprintf("%s:%s", cfg.App.AppHost, cfg.App.AppPort)
	logger.Info().Msgf("running REST server at '%s'", appAddr)
	if err := http.ListenAndServe(appAddr, router); err != nil {
		return fmt.Errorf("[NewRestApi] listen and serve: %w", err)
	}

	return nil
}
