package rest

import (
	"encoding/json"
	"errors"
	appErrors "gateway/internal/app_errors"
	"gateway/internal/models"
	"gateway/pkg/ctxutil"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strconv"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description This endpoint registers a new user in the system.
// @Tags users, v1
// @Accept json
// @Produce json
// @Param newUser body models.CreateUserDTO true "New User"
// @Success 200 {object} map[string]int "user_id"
// @Router /v1/users/register [post]
func (h *GatewayHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.RegisterUser")
	defer span.Finish()

	// Обработка запроса на регистрацию нового пользователя.
	var newUser = new(models.CreateUserDTO)
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[RegisterUser] unmarshal:%s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	userID, err := h.gatewayService.RegisterUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, appErrors.ErrUsernameOrEmailIsUsed) {
			h.ErrorUsernameOrEmailAlreadyUsed(w)
			return
		}

		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[RegisterUser] register:%s", err)
		h.ErrorInternalApi(w)
		return
	}

	// упаковываем данные для передачи пользователю
	response := struct {
		UserID int `json:"user_id"`
	}{UserID: userID}

	h.JSONSuccessRespond(w, response)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieves user details by their unique ID.
// @Tags users, v1
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserDTO
// @Router /v1/users/{id} [get]
func (h *GatewayHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.GetUserById")
	defer span.Finish()

	// Обработка запроса на получение информации о пользователе.
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetUserById] get id from query:%s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	user, err := h.gatewayService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			h.ErrorNotFound(w)
			return
		}

		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[GetUserById] get user:%s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем данные пользователю
	h.JSONSuccessRespond(w, user)
}

// UpdateUser godoc
// @Summary Update user information
// @Description Updates the information of an existing user.
// @Tags users, v1
// @Accept json
// @Produce json
// @Param updatedUser body models.UserDTO true "Updated User"
// @Success 200
// @Router /v1/users/update [put]
func (h *GatewayHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.UpdateUser")
	defer span.Finish()

	// Обработка запроса на обновление информации о пользователе.
	var updatedUser = new(models.UserDTO)
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdateUser] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	response, err := h.gatewayService.UpdateUser(ctx, updatedUser)
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdateUser] update user: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// передаем ответ пользователю
	h.JSONSuccessRespond(w, response)
}

// UpdatePassword godoc
// @Summary Update user's password
// @Description Allows a user to update their password.
// @Tags users, v1
// @Accept json
// @Produce json
// @Param passwordRequest body models.UpdateUserPasswordDTO true "Password Update Request"
// @Success 200
// @Router /v1/users/update-password [put]
func (h *GatewayHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.UpdatePassword")
	defer span.Finish()

	// Обработка запроса на изменение пароля пользователя.
	var passwordRequest = new(models.UpdateUserPasswordDTO)
	if err := json.NewDecoder(r.Body).Decode(&passwordRequest); err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdatePassword] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	if err := h.gatewayService.UpdatePassword(ctx, passwordRequest); err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UpdatePassword] update password: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ - в данном случе просто status 200
	h.JSONSuccessRespond(w, nil)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a user from the system based on their ID.
// @Tags users, v1
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200
// @Router /v1/users/delete/{id} [delete]
func (h *GatewayHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.DeleteUser")
	defer span.Finish()

	// Обработка запроса на удаление пользователя.
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[DeleteUser] get id from query: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	if err := h.gatewayService.DeleteUser(ctx, userID); err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[DeleteUser] delete user: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ - в данном случе просто status 200
	h.JSONSuccessRespond(w, nil)
}

// UserLogin godoc
// @Summary User login
// @Description Authenticates a user and returns access and refresh tokens.
// @Tags users, v1
// @Accept json
// @Produce json
// @Param request body models.UserLoginDTO true "Login Credentials"
// @Success 200 {object} map[string]string "Token Info"
// @Router /v1/users/login [post]
func (h *GatewayHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.UserLogin")
	defer span.Finish()

	// Обработка запроса на удаление пользователя.
	var request = new(models.UserLoginDTO)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UserLogin] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	response, err := h.gatewayService.Login(ctx, request)
	if err != nil {
		if errors.Is(err, appErrors.ErrWrongCredentials) {
			h.ErrorWrongCredentials(w)
			return
		}

		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[UserLogin] login: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ
	h.JSONSuccessRespond(w, response)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refreshes the user's access token.
// @Tags users, v1
// @Accept json
// @Produce json
// @Param request body models.UserTokens true "Token Refresh Request"
// @Success 200 {object} map[string]string "New Token Info"
// @Router /v1/users/refresh [post]
func (h *GatewayHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId, _ := ctxutil.GetRequestIDFromContext(ctx)
	span, ctx := opentracing.StartSpanFromContext(ctx, "gateway.Refresh")
	defer span.Finish()

	// Обработка запроса на удаление пользователя.
	var request = new(models.UserTokens)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[Refresh] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные слою бизнес-логики
	response, err := h.gatewayService.Refresh(ctx, request.RefreshToken)
	if err != nil {
		h.logger.Error().
			Str("requestId", requestId).
			Msgf("[Refresh] refresh: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ
	h.JSONSuccessRespond(w, response)
}
