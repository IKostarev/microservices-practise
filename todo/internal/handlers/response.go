package handlers

import (
	"encoding/json"
	"fmt"
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

func (h *TodoHandler) JSONErrorRespond(w http.ResponseWriter, httpCode int, errString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	if _, err := w.Write([]byte(errString)); err != nil {
		_ = fmt.Errorf("[JSONErrorRespond] error is - %w\n", err)
	}
}

func (h *TodoHandler) JSONSuccessRespond(w http.ResponseWriter, httpCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		w.WriteHeader(httpCode)
		return
	}

	jsonMarshal, err := json.Marshal(data)
	if err != nil {
		h.JSONErrorRespond(w, InternalServerError, "[JSONSuccessRespond] have error is marshal")
		return
	}

	w.WriteHeader(httpCode)

	if _, errWriter := w.Write(jsonMarshal); errWriter != nil {
		_ = fmt.Errorf("[JSONErrorRespond] error is - %w\n", errWriter)
	}
}
