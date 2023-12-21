package rest

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"todo/internal/api"
	"todo/internal/models"
)

type TodoHandler struct {
	todoService api.TodoService
	logger      *zerolog.Logger
}

func NewTodoHandler(todoService api.TodoService, logger *zerolog.Logger) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
		logger:      logger,
	}
}

func (h *TodoHandler) CreateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newTodo = new(models.TodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		h.logger.Err(err).Msg("[CreateToDoHandler] error unmarshal")
		h.JSONErrorRespond(w, BadRequest)
		return
	}

	todoID, err := h.todoService.CreateToDo(ctx, (*models.TodoDAO)(newTodo))
	if err != nil {
		h.logger.Err(err).Msg("[CreateToDoHandler] error create")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, Created, todoID)
}

func (h *TodoHandler) GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Err(err).Msg("[GetToDoHandler] error parse uuid")
		h.JSONErrorRespond(w, BadRequest)
		return
	}

	todo, err := h.todoService.GetToDo(ctx, todoID)
	if err != nil {
		h.logger.Err(err).Msg("[GetToDoHandler] error get todo")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, todo)
}

func (h *TodoHandler) GetToDosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID := uuid.Must(uuid.FromBytes([]byte(mux.Vars(r)["id"])))

	todos, err := h.todoService.GetToDos(ctx, todoID)
	if err != nil {
		h.logger.Err(err).Msg("[GetToDosHandler] error get ToDos")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, todos)
}

func (h *TodoHandler) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updTodo = new(models.TodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&updTodo); err != nil {
		h.logger.Err(err).Msg("[UpdateToDoHandler] error unmarshal")
		h.JSONErrorRespond(w, BadRequest)
		return
	}

	err := h.todoService.UpdateToDo(ctx, (*models.TodoDAO)(updTodo))
	if err != nil {
		h.logger.Err(err).Msg("[UpdateToDoHandler] error update todo")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, nil)
}

func (h *TodoHandler) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todoID := uuid.Must(uuid.FromBytes([]byte(mux.Vars(r)["id"])))

	if err := h.todoService.DeleteToDo(ctx, todoID); err != nil {
		h.logger.Err(err).Msg("[DeleteToDoHandler] error delete todo")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	h.JSONSuccessRespond(w, OK, nil)
}
