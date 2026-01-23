package utils

import (
	"encoding/json"
	"net/http"
)

//! Response structure générique
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

//! JSON envoie une réponse JSON
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

//! Success envoie une réponse de succès
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	JSON(w, statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

//! Error envoie une réponse d'erreur
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, Response{
		Success: false,
		Error:   message,
	})
}

//! Created - 201
func Created(w http.ResponseWriter, message string, data interface{}) {
	Success(w, http.StatusCreated, message, data)
}

//! OK - 200
func OK(w http.ResponseWriter, message string, data interface{}) {
	Success(w, http.StatusOK, message, data)
}

//! BadRequest - 400
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

//! Unauthorized - 401
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

//! Forbidden - 403
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

//! NotFound - 404
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

//! Conflict - 409
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, message)
}

//! InternalServerError - 500
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}
