package routes

import (
	"database/sql"
	"educnet/internal/auth"
	"educnet/internal/handler"
	"educnet/internal/middleware"
	"educnet/internal/repository"
	"educnet/internal/usecase"

	"github.com/gorilla/mux"
)

type Handlers struct {
	School  *handler.SchoolHandler
	Teacher *handler.TeacherHandler
	Student *handler.StudentHandler
	Auth    *handler.AuthHandler
	User    *handler.UserHandler
	Admin   *handler.AdminHandler
	Profile *handler.ProfileHandler
	Class   *handler.ClassHandler
	Subject *handler.SubjectHandler
	Chat    *handler.ChatHandler
}

func NewRouter(
	db *sql.DB,
	jwtService *auth.JWTService,
	jwtSecret string,
	//! REPOSITORIES
	schoolRepo repository.SchoolRepository,
	userRepo repository.UserRepository,
	subjectRepo repository.SubjectRepository,
	classRepo repository.ClassRepository,
	teacherSubjectRepo repository.TeacherSubjectRepository,
	studentClassRepo repository.StudentClassRepository,
	messageRepository repository.MessageRepository,
) *mux.Router {

	//! ========== USECASES ==========
	schoolUseCase := usecase.NewSchoolUseCase(db, schoolRepo, userRepo, jwtSecret) // ✅ FIXÉ
	teacherUseCase := usecase.NewTeacherUseCase(db, userRepo, schoolRepo, subjectRepo, teacherSubjectRepo)
	studentUseCase := usecase.NewStudentUseCase(db, userRepo, schoolRepo, classRepo, studentClassRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtService)
	adminUseCase := usecase.NewAdminUseCase(userRepo, teacherSubjectRepo, studentClassRepo, subjectRepo, classRepo)
	profileUseCase := usecase.NewProfileUseCase(userRepo, subjectRepo, classRepo, teacherSubjectRepo, studentClassRepo, schoolRepo)
	classUsecase := usecase.NewClassUsecase(classRepo)
	subjectUsecase := usecase.NewSubjectUsecase(subjectRepo)
	messageUsecase := usecase.NewMessageUseCase(messageRepository)
	//! ========== HANDLERS ==========
	handlers := &Handlers{
		School:  handler.NewSchoolHandler(schoolUseCase),
		Teacher: handler.NewTeacherHandler(teacherUseCase),
		Student: handler.NewStudentHandler(studentUseCase),
		Auth:    handler.NewAuthHandler(authUseCase),
		User:    handler.NewUserHandler(userRepo),
		Admin:   handler.NewAdminHandler(adminUseCase),
		Profile: handler.NewProfileHandler(profileUseCase),
		Class:   handler.NewClassHandler(classUsecase),
		Subject: handler.NewSubjectHandler(subjectUsecase),
		Chat:    handler.NewChatHandler(messageUsecase),
	}

	r := mux.NewRouter()

	r.Use(middleware.CORS)
	r.Use(middleware.Logger)

	api := r.PathPrefix("/api").Subrouter()

	//! ========== SUB-ROUTERS ==========
	SetupPublicRoutes(api, handlers)
	SetupUserRoutes(api, handlers, jwtService)
	SetupAdminRoutes(api, handlers, jwtService)
	SetupProfileRoutes(api, handlers, jwtService)
	SetupTeacherRoutes(api, handlers, jwtService)
	SetupStudentRoutes(api, handlers, jwtService)
	SetupWebSocketRoutes(api, handlers, jwtService)

	return r
}
