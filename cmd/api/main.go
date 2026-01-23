package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"educnet/internal/config"
	"educnet/internal/db"
	"educnet/internal/utils"
	"educnet/internal/middleware"
)

func main() {
	//! 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	//! 2. Connect to database
	database, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close(database)

	//! 3. Setup router
	r := mux.NewRouter()

	//! 4. Middleware
	r.Use(middleware.CORS)
	r.Use(middleware.Logging)

	//! 5. Routes
	api := r.PathPrefix("/api").Subrouter()

	//! Health check
	api.HandleFunc("/health", healthCheckHandler).Methods("GET")

	//! Hello World (Test de utils.Response)
	api.HandleFunc("/hello", helloWorldHandler).Methods("GET")

	//! 6. Start server
	addr := ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on http://localhost%s (env: %s)\n", addr, cfg.Server.Env)
	log.Printf("📍 Health check: http://localhost%s/api/health\n", addr)
	log.Printf("👋 Hello World: http://localhost%s/api/hello\n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

//! Handlers
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.OK(w, "Server is running", map[string]string{
		"status": "healthy",
		"env":    "development",
	})
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	utils.OK(w, "Hello from Educ-Net API!", map[string]interface{}{
		"message": "Welcome to Educ-Net Backend",
		"version": "1.0.0",
		"author":  "Your Name",
		"endpoints": []string{
			"GET /api/health",
			"GET /api/hello",
		},
	})
}
