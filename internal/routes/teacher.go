package routes

import (
	"educnet/internal/auth"
	"educnet/internal/middleware"

	"github.com/gorilla/mux"
)

// SetupTeacherRoutes configure les routes enseignant
func SetupTeacherRoutes(api *mux.Router, h *Handlers, jwtService *auth.JWTService) {
	// Teacher routes (JWT + TeacherOnly middleware)
	teacher := api.PathPrefix("/teacher").Subrouter()
	teacher.Use(middleware.JWTAuth(jwtService))
	teacher.Use(middleware.RoleRequired("teacher")) // À créer

	// ========== MY SUBJECTS ==========
	teacher.HandleFunc("/subjects", h.Teacher.GetMySubjects).Methods("GET")

	// ========== MY CLASSES ==========
	// teacher.HandleFunc("/classes", h.Teacher.GetMyClasses).Methods("GET")

	// ========== STUDENTS ==========
	// teacher.HandleFunc("/students", h.Teacher.GetMyStudents).Methods("GET")

	// ========== GRADES ==========
	// teacher.HandleFunc("/grades", h.Teacher.CreateGrade).Methods("POST")
	// teacher.HandleFunc("/grades/{id}", h.Teacher.UpdateGrade).Methods("PUT")

	// ========== ATTENDANCE ==========
	// teacher.HandleFunc("/attendance", h.Teacher.TakeAttendance).Methods("POST")
	// teacher.HandleFunc("/attendance", h.Teacher.GetAttendance).Methods("GET")
}
