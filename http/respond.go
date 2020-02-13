package http

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *handlers) respond(
	w http.ResponseWriter, _ *http.Request,
	msg string, data interface{}, status int,
) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response{
		Message: msg,
		Data:    data,
	})
	if err != nil {
		// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error
		h.logger.Printf("could not encode response to output: %v", err)
	}
}
