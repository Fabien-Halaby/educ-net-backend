package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"educnet/internal/config"
	"educnet/internal/db"
	"educnet/internal/handler"
	"educnet/internal/middleware"
	"educnet/internal/repository"
	"educnet/internal/usecase"
	"educnet/internal/utils"
	"educnet/internal/auth"
)

func main() {
	//! 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	//! 2. Connect to database
	database, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(database)

	//! Initialize JWT service
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)
	log.Println("JWT Secret: ", cfg.JWT.Secret)
	log.Printf("Access Token TTL: %d hours", cfg.JWT.AccessTokenTTL)

	//! 3. Initialize repositories (Data layer)
	schoolRepo := repository.NewSchoolRepository(database)
	userRepo := repository.NewUserRepository(database)
	subjectRepo := repository.NewSubjectRepository(database)
	classRepo := repository.NewClassRepository(database)
	teacherSubjectRepo := repository.NewTeacherSubjectRepository(database)
	studentClassRepo := repository.NewStudentClassRepository(database)

	//! 4. Initialize use cases (Business logic layer)
	schoolUseCase := usecase.NewSchoolUseCase(
		database,
		schoolRepo,
		userRepo,
		cfg.JWT.Secret,
	)
	teacherUseCase := usecase.NewTeacherUseCase(
		database,
		userRepo,
		schoolRepo,
		subjectRepo,
		teacherSubjectRepo,
	)
	studentUseCase := usecase.NewStudentUseCase(
		database,
		userRepo,
		schoolRepo,
		classRepo,
		studentClassRepo,
	)
	authUseCase := usecase.NewAuthUseCase(
		userRepo,
		jwtService,
	)

	//! 5. Initialize handlers (Presentation layer)
	schoolHandler := handler.NewSchoolHandler(schoolUseCase)
	teacherHandler := handler.NewTeacherHandler(teacherUseCase)
	studentHandler := handler.NewStudentHandler(studentUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userRepo)


	//! 6. Setup router
	r := mux.NewRouter()

	//! 7. Global middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logger)

	//! 8. Routes
	api := r.PathPrefix("/api").Subrouter()

	//! Health check
	api.HandleFunc("/health", health).Methods("GET")

	//! School routes
	api.HandleFunc("/schools/register", schoolHandler.CreateSchool).Methods("POST")

	//! Teacher routes
	api.HandleFunc("/teachers/register", teacherHandler.Register).Methods("POST")

	//! Student routes
	api.HandleFunc("/students/register", studentHandler.Register).Methods("POST")

	//! Auth routes
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	//! User routes (protected)
	userRouter := api.PathPrefix("/me").Subrouter()
	userRouter.Use(middleware.JWTAuth(jwtService))
	userRouter.HandleFunc("", userHandler.GetMe).Methods("GET")

	//! 9. Start server
	addr := ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on http://localhost%s (env: %s)\n", addr, cfg.Server.Env)
	log.Printf("📍 Health: http://localhost%s/api/health\n", addr)
	log.Printf("🏫 Register School: POST http://localhost%s/api/schools/register\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	utils.OK(w, "Server is running", map[string]string{
		"status": "healthy",
	})
}
