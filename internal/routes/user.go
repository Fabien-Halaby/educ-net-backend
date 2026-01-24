package routes

import (
	"educnet/internal/auth"
	"educnet/internal/middleware"

	"github.com/gorilla/mux"
)

// SetupUserRoutes configure les routes utilisateur (authentification requise)
func SetupUserRoutes(api *mux.Router, h *Handlers, jwtService *auth.JWTService) {
	// Protected routes (tous les utilisateurs authentifiés)
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTAuth(jwtService))

	// Profile
	protected.HandleFunc("/me", h.User.GetMe).Methods("GET")
	// protected.HandleFunc("/me", h.User.UpdateProfile).Methods("PUT")              // À venir
	// protected.HandleFunc("/me/password", h.User.ChangePassword).Methods("PUT")    // À venir
	// protected.HandleFunc("/me/avatar", h.User.UploadAvatar).Methods("POST")       // À venir

	// Subjects & Classes (à implémenter plus tard via AdminHandler)
	// protected.HandleFunc("/subjects", h.Subject.List).Methods("GET")
	// protected.HandleFunc("/subjects/{id}", h.Subject.GetByID).Methods("GET")
	// protected.HandleFunc("/classes", h.Class.List).Methods("GET")
	// protected.HandleFunc("/classes/{id}", h.Class.GetByID).Methods("GET")
}
