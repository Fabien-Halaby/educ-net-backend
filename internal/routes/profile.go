package routes

import (
	"educnet/internal/auth"
	"educnet/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// ! SetupProfileRoutes configure les routes de profil (protégées par JWT)
func SetupProfileRoutes(api *mux.Router, h *Handlers, jwtService *auth.JWTService) {
	//! Routes accessibles à tous (authentifiés)
	profile := api.PathPrefix("/me").Subrouter()
	profile.Use(middleware.JWTAuth(jwtService))

	profile.HandleFunc("", h.Profile.GetProfile).Methods("GET")
	profile.HandleFunc("", h.Profile.UpdateProfile).Methods("PUT")
	profile.HandleFunc("/password", h.Profile.ChangePassword).Methods("PUT")
	profile.HandleFunc("/avatar", h.Profile.UploadAvatar).Methods("POST")
	profile.HandleFunc("/school", h.Profile.GetSchool).Methods("GET")

	//! Routes ADMIN ONLY
	admin := profile.PathPrefix("/school").Subrouter()
	admin.Use(middleware.RoleRequired("admin"))

	admin.HandleFunc("", h.Profile.UpdateSchool).Methods("PUT")
	admin.HandleFunc("/logo", h.Profile.UploadSchoolLogo).Methods("POST")

	//! Routes teacher/student
	profile.HandleFunc("/subjects", h.Profile.GetMySubjects).Methods("GET")
	profile.HandleFunc("/class", h.Profile.GetMyClass).Methods("GET")

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
}
