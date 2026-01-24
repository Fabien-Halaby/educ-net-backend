package routes

import (
	"educnet/internal/auth"
	"educnet/internal/middleware"

	"github.com/gorilla/mux"
)

// SetupStudentRoutes configure les routes étudiant
func SetupStudentRoutes(api *mux.Router, h *Handlers, jwtService *auth.JWTService) {
	// Student routes (JWT + StudentOnly middleware)
	student := api.PathPrefix("/student").Subrouter()
	student.Use(middleware.JWTAuth(jwtService))
	student.Use(middleware.RoleRequired("student")) // À créer

	// ========== MY CLASS ==========
	// student.HandleFunc("/class", h.Student.GetMyClass).Methods("GET")
	
	// ========== MY SUBJECTS ==========
	// student.HandleFunc("/subjects", h.Student.GetMySubjects).Methods("GET")
	
	// ========== MY GRADES ==========
	// student.HandleFunc("/grades", h.Student.GetMyGrades).Methods("GET")
	
	// ========== MY ATTENDANCE ==========
	// student.HandleFunc("/attendance", h.Student.GetMyAttendance).Methods("GET")
	
	// ========== MY TEACHERS ==========
	// student.HandleFunc("/teachers", h.Student.GetMyTeachers).Methods("GET")
}
