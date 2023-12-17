package app

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"todo/config"
	"todo/internal/handlers"
	"todo/internal/repository"
	"todo/internal/service"
	"todo/pkg/db/postgres"
	"todo/pkg/logger"
)

type App struct {
	cfg    *config.Config
	router *mux.Router
	logger *zerolog.Logger
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger.NewLogger(cfg.Logger),
	}
}

func (a *App) RunAPI() error {
	conn, err := postgres.NewPostgres(&a.cfg.Database)
	if err != nil {
		a.logger.Err(err).Msgf("[RunAPI] db conn: %w\n", err)
		return err
	}
	defer conn.Close()

	repo := repository.NewTodoRepository(conn)
	serv := service.NewTodoService(repo)

	todoHandler := handlers.NewTodoHandler(serv, a.logger)

	a.router = mux.NewRouter()

	r := a.router.PathPrefix("/api/v1/todos").Subrouter()

	r.HandleFunc("/", todoHandler.CreateToDoHandler).Methods(http.MethodPost)
	r.HandleFunc("/batch", todoHandler.GetToDosHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", todoHandler.GetToDoHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", todoHandler.UpdateToDoHandler).Methods(http.MethodPut)
	r.HandleFunc("/{id}", todoHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	if err = http.ListenAndServe(a.cfg.App.AppPort, a.router); err != nil {
		a.logger.Err(err).Msgf("ListenAndServe error is - %s\n", err)
		return err
	}

	return nil
}
