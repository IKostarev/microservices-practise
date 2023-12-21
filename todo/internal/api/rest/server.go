package rest

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"todo/config"
	"todo/internal/api"
)

func NewRestAPI(cfg *config.Config, logger *zerolog.Logger, todoService api.TodoService) error {
	todoRestHandler := NewTodoHandler(todoService, logger)

	router := mux.NewRouter()

	router.HandleFunc("/todos/", todoRestHandler.CreateToDoHandler).Methods(http.MethodPost)
	router.HandleFunc("/todos/batch", todoRestHandler.GetToDosHandler).Methods(http.MethodGet)
	router.HandleFunc("/todos/{id}", todoRestHandler.GetToDoHandler).Methods(http.MethodGet)
	router.HandleFunc("/todos/{id}", todoRestHandler.UpdateToDoHandler).Methods(http.MethodPut)
	router.HandleFunc("/todos/{id}", todoRestHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(cfg.App.AppPort, router); err != nil {
		logger.Err(err).Msgf("ListenAndServe error is - %s\n", err)
		return err
	}

	return nil
}
