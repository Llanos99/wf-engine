package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Error writes an error response
func Error(w http.ResponseWriter, status int, err string, message string) {
	JSON(w, status, ErrorResponse{
		Error:   err,
		Message: message,
	})
}

// BadRequest writes a 400 error
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "bad_request", message)
}

// NotFound writes a 404 error
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "not_found", message)
}

// InternalError writes a 500 error
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, "internal_error", message)
}

// Conflict writes a 409 error
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "conflict", message)
}
