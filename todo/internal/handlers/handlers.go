package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"todo/internal/models"
)

type TodoHandler struct {
	todoService TodoService
}

func NewTodoHandler(todoService TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

func (h *TodoHandler) CreateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newTodo = new(models.TodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		fmt.Println("[CreateToDoHandler] error unmarshal")
		h.JSONErrorRespond(w, BadRequest)
		return
	}

	todoID, err := h.todoService.CreateToDo(ctx, newTodo)
	if err != nil {
		// TODO обработать все возможные ошибки
		fmt.Println("[CreateToDoHandler] error create")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	resp := struct {
		TodoID uuid.UUID `json:"todo_id"`
	}{TodoID: todoID.ID}

	h.JSONSuccessRespond(w, Created, resp)
}

func (h *TodoHandler) GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("[GetToDoHandler] error parse uuid")
		h.JSONErrorRespond(w, BadRequest)
		return
	}

	todo, err := h.todoService.GetToDo(ctx, todoID)
	if err != nil {
		// TODO обработать все возможные ошибки
		fmt.Println("[GetToDoHandler] error get todo")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, todo)
}

func (h *TodoHandler) GetToDosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todos, err := h.todoService.GetToDos(ctx)
	if err != nil {
		// TODO обработать все возможные ошибки
		fmt.Println("[GetToDosHandler] error get ToDos")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, todos)
}

func (h *TodoHandler) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updTodo = new(models.TodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&updTodo); err != nil {
		fmt.Println("[UpdateToDoHandler] error unmarshal")
		h.JSONErrorRespond(w, BadRequest)
		return
	}

	resp, err := h.todoService.UpdateToDo(ctx, updTodo)
	if err != nil {
		// TODO обработать все возможные ошибки
		fmt.Println("[UpdateToDoHandler] error update todo")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, resp)
}

func (h *TodoHandler) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID := uuid.Must(uuid.FromBytes([]byte(mux.Vars(r)["id"])))

	if err := h.todoService.DeleteToDo(ctx, todoID); err != nil {
		fmt.Println("[DeleteToDoHandler] error delete todo")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, nil)
}
