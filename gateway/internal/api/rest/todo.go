package rest

import (
	"encoding/json"
	"errors"
	"gateway/internal/app_errors"
	"gateway/internal/models"
	"gateway/pkg/ctxutil"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strconv"
)

// CreateToDoHandler godoc
// @Summary Create a new todo
// @Description This endpoint creates a new todo in the system.
// @Tags todo, v1
// @Accept json
// @Produce json
// @Param newTodo body models.CreateTodoDTO true "New Todo"
// @Success 200 {object} map[string]int "todo_id"
// @Router /v1/todos [post]
func (h *GatewayHandler) CreateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.CreateToDo")
	defer span.Finish()

	var newTodo = new(models.CreateTodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[CreateToDoHandler] unmarshal: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	createdTodo, err := h.gatewayService.CreateToDo(ctx, newTodo)
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[CreateToDoHandler] create todo: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, createdTodo)
}

// GetToDoHandler godoc
// @Summary Get todo by ID
// @Description Retrieves todo details by their unique ID.
// @Tags todo, v1
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.TodoDTO
// @Route /v1/todos/{id} [get]
func (h *GatewayHandler) GetToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.GetTodo")
	defer span.Finish()

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetToDoHandler] parse id from url: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	todo, err := h.gatewayService.GetToDo(ctx, id)
	if err != nil {
		if errors.Is(err, app_errors.ErrNotFound) {
			h.ErrorNotFound(w)
			return
		}

		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetToDoHandler] get todo: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, todo)
}

// GetToDosHandler godoc
// @Summary Get todos by ID
// @Description Retrieves todos details by their unique ID.
// @Tags todo, v1
// @Accept json
// @Produce json
// @Param id path int true "Todos ID"
// @Success 200 {object} models.TodoDTO
// @Route /v1/todos/batch [get]
func (h *GatewayHandler) GetToDosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.GetTodos")
	defer span.Finish()

	getTodos, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetToDosHandler] parse id from url: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	todos, err := h.gatewayService.GetToDos(ctx, getTodos)
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetToDosHandler] get todos: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, todos)
}

// UpdateToDoHandler godoc
// @Summary Update todo information
// @Description Updates the information of an existing todo.
// @Tags todo, v1
// @Accept json
// @Produce json
// @Param updateTodo body models.UpdateTodoDTO true "Update ToDo"
// @Success 200
// @Router /v1/todos/{id} [put]
func (h *GatewayHandler) UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.UpdateToDo")
	defer span.Finish()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdateToDoHandler] parse id from url: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	todo := new(models.UpdateTodoDTO)
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		h.ErrorBadRequest(w)
		return
	}
	todo.ID = id

	updatedTodo, err := h.gatewayService.UpdateToDo(ctx, todo)
	if err != nil {
		if errors.Is(err, app_errors.ErrNotFound) {
			h.ErrorNotFound(w)
			return
		}

		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdateToDoHandler] update todo: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, updatedTodo)
}

// DeleteToDoHandler godoc
// @Summary Delete a todo
// @Description Deletes a todo from the system based on their ID.
// @Tags todo, v1
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200
// @Router /v1/todos/{id} [delete]
func (h *GatewayHandler) DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.DeleteTodo")
	defer span.Finish()

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[DeleteToDoHandler] parse id from url: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	err = h.gatewayService.DeleteToDo(ctx, id)
	if err != nil {
		if errors.Is(err, app_errors.ErrNotFound) {
			h.ErrorNotFound(w)
			return
		}

		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[DeleteToDoHandler] delete todo: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	h.JSONSuccessRespond(w, nil)
}
