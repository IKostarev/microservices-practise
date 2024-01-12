package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
	"todo/internal/api"
	"todo/internal/models"
)

type TodoHandler struct {
	logger      *zerolog.Logger
	todoService api.TodoService
}

func NewTodoHandler(
	logger *zerolog.Logger,
	todoService api.TodoService,
) *TodoHandler {
	return &TodoHandler{
		logger:      logger,
		todoService: todoService,
	}
}

func (h *TodoHandler) CreateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newTodo = new(models.CreateTodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		h.logger.Err(err).Msg("[CreateToDoHandler] error unmarshal")
		h.ErrorBadRequest(w)
		return
	}

	todoID, err := h.todoService.CreateToDo(ctx, newTodo)
	if err != nil {
		h.logger.Err(err).Msg("[CreateToDoHandler] error create")
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, todoID)
}

func (h *TodoHandler) GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Err(err).Msg("[GetToDoHandler] error parse uuid")
		h.ErrorBadRequest(w)
		return
	}

	todo, err := h.todoService.GetToDo(ctx, todoID)
	if err != nil {
		h.logger.Err(err).Msg("[GetToDoHandler] error get todo")
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, todo)
}

func (h *TodoHandler) GetToDosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Err(err).Msg("[GetToDosHandler] error parse uuid")
		h.ErrorBadRequest(w)
		return
	}

	todos, err := h.todoService.GetToDos(ctx, todoID)
	if err != nil {
		h.logger.Err(err).Msg("[GetToDosHandler] error get ToDos")
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, todos)
}

func (h *TodoHandler) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updTodo = new(models.UpdateTodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&updTodo); err != nil {
		h.logger.Err(err).Msg("[UpdateToDoHandler] error unmarshal")
		h.ErrorBadRequest(w)
		return
	}

	todoID, err := h.todoService.UpdateToDo(ctx, updTodo)
	if err != nil {
		h.logger.Err(err).Msg("[UpdateToDoHandler] error update todo")
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, todoID)
}

func (h *TodoHandler) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Err(err).Msg("[GetToDosHandler] error parse uuid")
		h.ErrorBadRequest(w)
		return
	}

	if err := h.todoService.DeleteToDo(ctx, todoID); err != nil {
		h.logger.Err(err).Msg("[DeleteToDoHandler] error delete todo")
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, todoID)
}
