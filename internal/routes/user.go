package routes

import (
	"educnet/internal/auth"
	"educnet/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(api *mux.Router, h *Handlers, jwtService *auth.JWTService) {
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTAuth(jwtService))

	//! Profile management
	// protected.HandleFunc("/me", h.Profile.GetProfile).Methods("GET")
	// protected.HandleFunc("/me", h.Profile.UpdateProfile).Methods("PUT")
	// protected.HandleFunc("/me/password", h.Profile.ChangePassword).Methods("PUT")

	//! Role-specific
	protected.HandleFunc("/me/subjects", h.Profile.GetMySubjects).Methods("GET") // Teachers
	protected.HandleFunc("/me/classes", h.Profile.GetMyClass).Methods("GET")     // Students
	// Subjects & Classes (à implémenter plus tard via AdminHandler)
	// protected.HandleFunc("/subjects", h.Subject.List).Methods("GET")
	// protected.HandleFunc("/subjects/{id}", h.Subject.GetByID).Methods("GET")
	// protected.HandleFunc("/classes", h.Class.List).Methods("GET")
	// protected.HandleFunc("/classes/{id}", h.Class.GetByID).Methods("GET")
}
