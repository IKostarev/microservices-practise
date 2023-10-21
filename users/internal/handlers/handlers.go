package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
	"users/internal/models"
)

type UserHandler struct {
	logger      *zerolog.Logger
	userService UserService
}

func NewUserHandler(
	logger *zerolog.Logger,
	userService UserService,
) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Обработка запроса на регистрацию нового пользователя.
	var newUser = new(models.CreateUserDTO)
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.logger.Error().Msgf("[RegisterUser] unmarshal:%s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	userID, err := h.userService.RegisterUser(ctx, newUser)
	if err != nil {
		h.logger.Error().Msgf("[RegisterUser] register:%s", err)
		h.ErrorInternalApi(w)
		return
	}

	// упаковываем данные для передачи пользователю
	response := struct {
		UserID int `json:"user_id"`
	}{UserID: userID}

	h.JSONSuccessRespond(w, response)
}

func (h *UserHandler) GetGetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка запроса на получение информации о пользователе.
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().Msgf("[GetGetUserById] get id from query:%s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		h.logger.Error().Msgf("[GetGetUserById] get user:%s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем данные пользователю
	h.JSONSuccessRespond(w, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка запроса на обновление информации о пользователе.
	var updatedUser = new(models.UserDTO)
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		h.logger.Error().Msgf("[UpdateUser] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	response, err := h.userService.UpdateUser(ctx, updatedUser)
	if err != nil {
		h.logger.Error().Msgf("[UpdateUser] update user: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// передаем ответ пользователю
	h.JSONSuccessRespond(w, response)
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка запроса на изменение пароля пользователя.
	var passwordRequest = new(models.UpdateUserPasswordDTO)
	if err := json.NewDecoder(r.Body).Decode(&passwordRequest); err != nil {
		h.logger.Error().Msgf("[UpdatePassword] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	if err := h.userService.UpdatePassword(ctx, passwordRequest); err != nil {
		h.logger.Error().Msgf("[UpdatePassword] update password: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ - в данном случе просто status 200
	h.JSONSuccessRespond(w, nil)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка запроса на удаление пользователя.
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.logger.Error().Msgf("[DeleteUser] get id from query: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные в слой сервиса
	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		h.logger.Error().Msgf("[DeleteUser] delete user: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ - в данном случе просто status 200
	h.JSONSuccessRespond(w, nil)
}

func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка запроса на удаление пользователя.
	var request = new(models.UserLoginDTO)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error().Msgf("[UserLogin] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	response, err := h.userService.Login(ctx, request)
	if err != nil {
		h.logger.Error().Msgf("[UserLogin] login: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ
	h.JSONSuccessRespond(w, response)
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка запроса на удаление пользователя.
	var request = new(models.UserTokens)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error().Msgf("[Refresh] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем данные слою бизнес-логики
	response, err := h.userService.Refresh(ctx, request.RefreshToken)
	if err != nil {
		h.logger.Error().Msgf("[Refresh] refresh: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	// возвращаем пользователю ответ
	h.JSONSuccessRespond(w, response)
}

func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Обработка запроса на удаление пользователя.
	var request = new(models.UserTokens)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error().Msgf("[Verify] unmarshall: %s", err)
		h.ErrorBadRequest(w)
		return
	}

	// передаем токен на проверку в слой бизнес-логики
	userID, err := h.userService.VerifyToken(ctx, request.AccessToken)
	if err != nil {
		h.logger.Error().Msgf("[Verify] verify: %s", err)
		h.ErrorInternalApi(w)
		return
	}

	response := struct {
		UserID int `json:"user_id"`
	}{
		UserID: userID,
	}

	// возвращаем пользователю ответ
	h.JSONSuccessRespond(w, response)
}
