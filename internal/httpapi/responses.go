package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/baditaflorin/civitas/internal/storage"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error, fallback string) {
	status := http.StatusInternalServerError
	code := "internal_error"
	message := fallback
	if errors.Is(err, storage.ErrNotFound) {
		status = http.StatusNotFound
		code = "not_found"
		message = "resource not found"
	}
	writeJSON(w, status, errorResponse{Code: code, Message: message})
}

func decodeJSON(r *http.Request, value any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(value)
}
