package main

import (
	"log"
	"net/http"

	"educnet/internal/auth"
	"educnet/internal/config"
	"educnet/internal/db"
	"educnet/internal/middleware"
	"educnet/internal/repository"
	"educnet/internal/routes"
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
	log.Println("✅ Database connected successfully")

	//! 3. Initialize JWT service
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)
	log.Printf("🔑 JWT configured (TTL: %dh)", cfg.JWT.AccessTokenTTL)

	//! 4. Initialize repositories
	schoolRepo := repository.NewSchoolRepository(database)
	userRepo := repository.NewUserRepository(database)
	subjectRepo := repository.NewSubjectRepository(database)
	classRepo := repository.NewClassRepository(database)
	teacherSubjectRepo := repository.NewTeacherSubjectRepository(database)
	studentClassRepo := repository.NewStudentClassRepository(database)
	messageRepository := repository.NewMessageRepository(database)
	//! 5. Setup router (all routes configured in routes package)
	router := routes.NewRouter(
		database,
		jwtService,
		cfg.JWT.Secret,
		schoolRepo,
		userRepo,
		subjectRepo,
		classRepo,
		teacherSubjectRepo,
		studentClassRepo,
		messageRepository,
	)

	handler := middleware.CORS(router)

	//! 6. Start server
	addr := ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on http://localhost%s (env: %s)", addr, cfg.Server.Env)
	log.Printf("📍 Health: http://localhost%s/api/health", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
