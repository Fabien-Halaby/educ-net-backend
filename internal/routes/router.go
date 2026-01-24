package routes

import (
	"educnet/internal/auth"
	"educnet/internal/handler"
	"educnet/internal/middleware"
	"educnet/internal/repository"
	"educnet/internal/usecase"
	"database/sql"

	"github.com/gorilla/mux"
)

// Handlers struct pour passer aux sous-routeurs
type Handlers struct {
	School  *handler.SchoolHandler
	Teacher *handler.TeacherHandler
	Student *handler.StudentHandler
	Auth    *handler.AuthHandler
	User    *handler.UserHandler
	Admin   *handler.AdminHandler
}

// SetupRouter configure toutes les routes de l'application
func SetupRouter(
	db *sql.DB,
	jwtService *auth.JWTService,
	jwtSecret string, // ✅ AJOUTÉ
	// Repositories
	schoolRepo repository.SchoolRepository,
	userRepo repository.UserRepository,
	subjectRepo repository.SubjectRepository,
	classRepo repository.ClassRepository,
	teacherSubjectRepo repository.TeacherSubjectRepository,
	studentClassRepo repository.StudentClassRepository,
) *mux.Router {

	// ========== INITIALIZE USE CASES ==========
	schoolUseCase := usecase.NewSchoolUseCase(db, schoolRepo, userRepo, jwtSecret) // ✅ FIXÉ
	teacherUseCase := usecase.NewTeacherUseCase(db, userRepo, schoolRepo, subjectRepo, teacherSubjectRepo)
	studentUseCase := usecase.NewStudentUseCase(db, userRepo, schoolRepo, classRepo, studentClassRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtService)
	adminUseCase := usecase.NewAdminUseCase(userRepo, teacherSubjectRepo, studentClassRepo, subjectRepo, classRepo)

	// ========== INITIALIZE HANDLERS ==========
	handlers := &Handlers{
		School:  handler.NewSchoolHandler(schoolUseCase),
		Teacher: handler.NewTeacherHandler(teacherUseCase),
		Student: handler.NewStudentHandler(studentUseCase),
		Auth:    handler.NewAuthHandler(authUseCase),
		User:    handler.NewUserHandler(userRepo),
		Admin:   handler.NewAdminHandler(adminUseCase),
	}

	// ========== SETUP ROUTER ==========
	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logger)

	// API prefix
	api := r.PathPrefix("/api").Subrouter()

	// ========== SETUP SUB-ROUTERS ==========
	SetupPublicRoutes(api, handlers)
	SetupUserRoutes(api, handlers, jwtService)
	SetupAdminRoutes(api, handlers, jwtService)
	// SetupTeacherRoutes(api, handlers, jwtService)  // À venir
	// SetupStudentRoutes(api, handlers, jwtService)  // À venir

	return r
}
