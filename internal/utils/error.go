package utils

import (
	"educnet/internal/domain"
	"errors"
	"log"
	"net/http"
)

func HandleUseCaseError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, `{"error":"Resource not found"}`, http.StatusNotFound)
	case errors.Is(err, domain.ErrForbidden):
		http.Error(w, `{"error":"Access forbidden"}`, http.StatusForbidden)
	case errors.Is(err, domain.ErrUnauthorized):
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
	case errors.Is(err, domain.ErrValidation):
		http.Error(w, `{"error":"Validation failed"}`, http.StatusUnprocessableEntity)
	default:
		log.Printf("Internal error: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
	}
}
