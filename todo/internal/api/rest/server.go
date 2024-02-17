package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/pprof"
	"todo/config"
	"todo/internal/api"
)

func NewRestApi(cfg *config.Config, logger *zerolog.Logger, todoService api.TodoService) error {
	todoRestHandler := NewTodoHandler(logger, todoService)

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
