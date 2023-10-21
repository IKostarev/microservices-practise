package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"todo/internal/handlers"
	"todo/internal/service"
)

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

	// Настройка роутеров и запуск сервера.
	a.router = mux.NewRouter()
	a.router.HandleFunc("/todos", todoHandler.CreateToDoHandler).Methods(http.MethodPost)
	a.router.HandleFunc("/todos/{id}", todoHandler.GetToDoHandler).Methods(http.MethodGet)
	a.router.HandleFunc("/todos/batch", todoHandler.GetToDosHandler).Methods(http.MethodGet)
	a.router.HandleFunc("/todos/{id}", todoHandler.UpdateToDoHandler).Methods(http.MethodPut)
	a.router.HandleFunc("/todos/{id}", todoHandler.DeleteToDoHandler).Methods(http.MethodDelete)

	http.ListenAndServe(":8080", a.router)
}
