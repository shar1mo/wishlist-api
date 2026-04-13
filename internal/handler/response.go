package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: message,
	})
}

func writeValidationError(w http.ResponseWriter, details map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, ValidationErrorResponse{
		Error:   "validation failed",
		Details: details,
	})
}