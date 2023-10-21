package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

func (h *UserHandler) ErrorBadRequest(w http.ResponseWriter) {
	h.JSONErrorRespond(w, http.StatusBadRequest, ErrBadRequest)
}

func (h *UserHandler) ErrorInternalApi(w http.ResponseWriter) {
	h.JSONErrorRespond(w, http.StatusInternalServerError, ErrInternalApi)
}

func (h *UserHandler) JSONErrorRespond(w http.ResponseWriter, httpCode int, err *ApiError) {
	// установка хэдера ответа
	w.Header().Set("Content-Type", "application/json")

	if err == (*ApiError)(nil) {
		w.WriteHeader(httpCode)
		return
	}

	data := ErrorResponse{
		ErrorCode:    string(err.ErrCode),
		ErrorMessage: err.Error(),
	}

	rawData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		h.JSONErrorRespond(w, http.StatusInternalServerError, NewApiError("can't encode response to json", ErrCodeInvalidJsonFormat))
	}

	w.WriteHeader(httpCode)

	_, _ = w.Write(rawData)
}

func (h *UserHandler) JSONSuccessRespond(w http.ResponseWriter, data interface{}) {
	// установка хэдера ответа
	w.Header().Set("Content-Type", "application/json")

	// если нет тела ответа, то просто вернем ответ со статусом 200
	if data == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	// если есть тело ответа, замарашалим его
	rawData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		// если возникла ошибка, связанная с маршаллингом, вернем код 500
		h.JSONErrorRespond(w, http.StatusInternalServerError, NewApiError(marshalErr.Error(), ErrCodeInvalidJsonFormat))
		return
	}

	// установим статус 200
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(rawData)
}
