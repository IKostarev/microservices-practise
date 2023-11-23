package app

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"todo/internal/handlers"
	"todo/internal/service"
)

// TODO перенести в конфиг
const PORT = ":8080"

type App struct {
	router      *mux.Router
	todoService handlers.TodoService
}

func NewApp() (*App, error) {
	todoService := service.NewTodoService(nil)

	return &App{
		todoService: todoService,
	}, nil
}

func (a *App) RunAPI() {
	todoHandler := handlers.NewTodoHandler(a.todoService)

	a.router = mux.NewRouter()

	r := a.router.PathPrefix("/api/v1/todos").Subrouter()

	r.HandleFunc("/", todoHandler.CreateToDoHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}", todoHandler.GetToDoHandler).Methods(http.MethodGet)
	r.HandleFunc("/batch", todoHandler.GetToDosHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}", todoHandler.UpdateToDoHandler).Methods(http.MethodPut)
	r.HandleFunc("/{id}", todoHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	if err := http.ListenAndServe(PORT, a.router); err != nil {
		log.Fatalf("ListenAndServe error is - %w\n", err)
	}
}
