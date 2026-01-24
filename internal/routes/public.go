package routes

import (
	"educnet/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
)

//! SetupPublicRoutes configure les routes publiques (sans authentification)
func SetupPublicRoutes(api *mux.Router, h *Handlers) {
	//! Health check
	api.HandleFunc("/health", health).Methods("GET")

	//! Registration
	api.HandleFunc("/schools/register", h.School.CreateSchool).Methods("POST")
	api.HandleFunc("/teachers/register", h.Teacher.Register).Methods("POST")
	api.HandleFunc("/students/register", h.Student.Register).Methods("POST")

	//! Authentication
	api.HandleFunc("/auth/login", h.Auth.Login).Methods("POST")
	// api.HandleFunc("/auth/refresh", h.Auth.RefreshToken).Methods("POST")  // À venir
	// api.HandleFunc("/auth/logout", h.Auth.Logout).Methods("POST")         // À venir
}

//! health handler
func health(w http.ResponseWriter, r *http.Request) {
	utils.OK(w, "Server is running", map[string]string{
		"status":  "healthy",
		"version": "1.0.0",
	})
}
