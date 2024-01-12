package rest

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}

func (h *TodoHandler) ErrorBadRequest(w http.ResponseWriter) {
	h.JSONErrorRespond(w, http.StatusBadRequest, ErrBadRequest)
}

func (h *TodoHandler) ErrorNotFound(w http.ResponseWriter) {
	h.JSONErrorRespond(w, http.StatusNotFound, ErrNotFound)
}

func (h *TodoHandler) ErrorInternalApi(w http.ResponseWriter) {
	h.JSONErrorRespond(w, http.StatusInternalServerError, ErrInternalApi)
}

func (h *TodoHandler) JSONErrorRespond(w http.ResponseWriter, httpCode int, err *ApiError) {
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
		h.logger.Error().Msgf("[JSONErrorRespond] marshal:%s", err)
		h.JSONErrorRespond(w, http.StatusInternalServerError, NewApiError("marshal to json", ErrCodeInvalidJsonFormat))
	}

	w.WriteHeader(httpCode)

	_, writeErr := w.Write(rawData)
	if writeErr != nil {
		h.logger.Error().Msgf("[JSONErrorRespond] write response:%s", writeErr)
	}
}

func (h *TodoHandler) JSONSuccessRespond(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	rawData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		h.logger.Error().Msgf("[JSONSuccessRespond] marshal:%s", marshalErr)
		h.JSONErrorRespond(w, http.StatusInternalServerError, NewApiError(marshalErr.Error(), ErrCodeInvalidJsonFormat))
		return
	}

	w.WriteHeader(http.StatusOK)

	_, writeErr := w.Write(rawData)
	if writeErr != nil {
		h.logger.Error().Msgf("[JSONSuccessRespond] write response:%s", writeErr)
	}
}
