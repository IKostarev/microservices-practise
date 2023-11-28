package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"todo/config"
	"todo/internal/handlers"
	"todo/internal/service"
)

type App struct {
	cfg         *config.Config
	router      *mux.Router
	todoService handlers.TodoService
}

func NewApp(cfg *config.Config) (*App, error) {
	todoService := service.NewTodoService(nil)

	return &App{
		cfg:         cfg,
		todoService: todoService,
	}, nil
}

func (a *App) RunAPI() {
	todoHandler := handlers.NewTodoHandler(a.todoService)

	a.router = mux.NewRouter()

	r := a.router.PathPrefix("/api/v1/todos").Subrouter()

	r.HandleFunc("/", todoHandler.CreateToDoHandler).Methods(http.MethodPost)
	r.HandleFunc("/batch", todoHandler.GetToDosHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", todoHandler.GetToDoHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", todoHandler.UpdateToDoHandler).Methods(http.MethodPut)
	r.HandleFunc("/{id}", todoHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(a.cfg.App.AppPort, a.router); err != nil {
		fmt.Printf("ListenAndServe error is - %s\n", err)
	}
}
