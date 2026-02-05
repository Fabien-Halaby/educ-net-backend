package routes

import (
	"educnet/internal/auth"
	"educnet/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupWebSocketRoutes(r *mux.Router, h *Handlers, jwtService *auth.JWTService) {
	wsRouter := r.PathPrefix("/ws").Subrouter()

	wsRouter.Use(middleware.JWTAuth(jwtService))

	wsRouter.HandleFunc("/chat/{classId}", h.Chat.HandleWebSocket).Methods("GET")

}
