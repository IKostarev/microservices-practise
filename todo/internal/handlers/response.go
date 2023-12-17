package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	OK      int = http.StatusOK
	Created int = http.StatusCreated
)

const (
	BadRequest          int = http.StatusBadRequest
	NotFound            int = http.StatusNotFound
	InternalServerError int = http.StatusInternalServerError
)

func (h *TodoHandler) JSONErrorRespond(w http.ResponseWriter, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
}

func (h *TodoHandler) JSONSuccessRespond(w http.ResponseWriter, httpCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		w.WriteHeader(httpCode)
		return
	}

	jsonMarshal, err := json.Marshal(data)
	if err != nil {
		h.logger.Err(err).Msg("[JSONSuccessRespond] have error is marshal")
		h.JSONErrorRespond(w, InternalServerError)
		return
	}

	w.WriteHeader(httpCode)

	if _, errWriter := w.Write(jsonMarshal); errWriter != nil {
		h.logger.Err(err).Msgf("[JSONErrorRespond] error is - %w\n", errWriter)
	}
}
