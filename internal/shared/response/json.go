package response

import (
	"encoding/json"
	"net/http"
)

const (
	statusSuccess = "success"
	statusError   = "error"
)

// Envelope is the standard JSON response wrapper.
type Envelope struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func write(w http.ResponseWriter, status int, env Envelope) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(env)
}

// Success writes a 200 OK response.
func Success(w http.ResponseWriter, data any) {
	write(w, http.StatusOK, Envelope{Status: statusSuccess, Message: "OK", Data: data})
}

// Created writes a 201 Created response.
func Created(w http.ResponseWriter, data any) {
	write(w, http.StatusCreated, Envelope{Status: statusSuccess, Message: "Created", Data: data})
}

// Error writes a JSON error response.
func Error(w http.ResponseWriter, status int, message string) {
	write(w, status, Envelope{Status: statusError, Message: message})
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
