package handlers

import "net/http"

type TodoHandler struct {
	todoService TodoService
}

func NewTodoHandler(todoService TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

func (h *TodoHandler) CreateToDoHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) GetToDoHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) GetToDosHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *TodoHandler) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {

}
