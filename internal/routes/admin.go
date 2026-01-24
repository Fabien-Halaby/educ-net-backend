package routes

import (
	"educnet/internal/auth"
	"educnet/internal/middleware"

	"github.com/gorilla/mux"
)

// SetupAdminRoutes configure les routes admin (authentification + rôle admin requis)
func SetupAdminRoutes(api *mux.Router, h *Handlers, jwtService *auth.JWTService) {
	// Admin routes (JWT + AdminOnly middleware)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.JWTAuth(jwtService))
	admin.Use(middleware.AdminOnly)

	// ========== USER MANAGEMENT ==========
	admin.HandleFunc("/users/pending", h.Admin.GetPendingUsers).Methods("GET")
	admin.HandleFunc("/users/{id}/approve", h.Admin.ApproveUser).Methods("POST")
	admin.HandleFunc("/users/{id}/reject", h.Admin.RejectUser).Methods("POST")
	admin.HandleFunc("/users", h.Admin.GetAllUsers).Methods("GET")
	// admin.HandleFunc("/users/{id}", h.Admin.GetUserByID).Methods("GET")           // À venir
	// admin.HandleFunc("/users/{id}/suspend", h.Admin.SuspendUser).Methods("POST")  // À venir

	// ========== SUBJECT MANAGEMENT (CRUD) - À IMPLÉMENTER ==========
	admin.HandleFunc("/subjects", h.Admin.CreateSubject).Methods("POST")
	admin.HandleFunc("/subjects/{id}", h.Admin.UpdateSubject).Methods("PUT")
	admin.HandleFunc("/subjects/{id}", h.Admin.DeleteSubject).Methods("DELETE")

	// ========== CLASS MANAGEMENT (CRUD) - À IMPLÉMENTER ==========
	admin.HandleFunc("/classes", h.Admin.CreateClass).Methods("POST")
	admin.HandleFunc("/classes/{id}", h.Admin.UpdateClass).Methods("PUT")
	admin.HandleFunc("/classes/{id}", h.Admin.DeleteClass).Methods("DELETE")

	// ========== DASHBOARD & STATS ==========
	// admin.HandleFunc("/dashboard", h.Admin.GetDashboard).Methods("GET")           // À venir
	// admin.HandleFunc("/stats", h.Admin.GetStats).Methods("GET")                   // À venir
}
