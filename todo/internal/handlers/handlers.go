package handlers

import (
	"encoding/json"
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
		h.JSONErrorRespond(w, BadRequest, "[CreateToDoHandler] error unmarshal")
		return
	}
	defer r.Body.Close()

	todoID, err := h.todoService.CreateToDo(ctx, newTodo)
	if err != nil {
		// TODO обработать все возможные ошибки
		h.JSONErrorRespond(w, InternalServerError, "[CreateToDoHandler] error create")
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
		h.JSONErrorRespond(w, BadRequest, "[GetToDoHandler] error parse uuid")
		return
	}

	todo, err := h.todoService.GetToDo(ctx, todoID)
	if err != nil {
		// TODO обработать все возможные ошибки
		h.JSONErrorRespond(w, InternalServerError, "[GetToDoHandler] error get todo")
		return
	}

	h.JSONSuccessRespond(w, OK, todo)
}

func (h *TodoHandler) GetToDosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todos, err := h.todoService.GetToDos(ctx)
	if err != nil {
		// TODO обработать все возможные ошибки
		h.JSONErrorRespond(w, InternalServerError, "[GetToDosHandler] error get ToDos")
		return
	}

	h.JSONSuccessRespond(w, OK, todos)
}

func (h *TodoHandler) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updTodo = new(models.TodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&updTodo); err != nil {
		h.JSONErrorRespond(w, BadRequest, "[UpdateToDoHandler] error unmarshal")
		return
	}
	defer r.Body.Close()

	resp, err := h.todoService.UpdateToDo(ctx, updTodo)
	if err != nil {
		// TODO обработать все возможные ошибки
		h.JSONErrorRespond(w, InternalServerError, "[UpdateToDoHandler] error update todo")
		return
	}

	h.JSONSuccessRespond(w, OK, resp)
}

func (h *TodoHandler) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID := uuid.Must(uuid.FromBytes([]byte(mux.Vars(r)["id"])))

	if err := h.todoService.DeleteToDo(ctx, todoID); err != nil {
		h.JSONErrorRespond(w, InternalServerError, "[DeleteToDoHandler] error delete todo")
		return
	}

	h.JSONSuccessRespond(w, OK, nil)
}
