package middleware

import (
	"context"
	"educnet/internal/auth"
	"educnet/internal/utils"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

// JWTAuth middleware pour protéger les routes
func JWTAuth(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			log.Printf("[AUTH] Authorization header: %s", authHeader) // DEBUG
			
			if authHeader == "" {
				log.Println("[AUTH] Missing authorization header") // DEBUG
				utils.Unauthorized(w, "Missing authorization header")
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Printf("[AUTH] Invalid format. Parts: %v", parts) // DEBUG
				utils.Unauthorized(w, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]
			log.Printf("[AUTH] Token: %s...", tokenString[:20]) // DEBUG (premiers 20 chars)

			// Validate token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				log.Printf("[AUTH] Token validation error: %v", err) // DEBUG
				utils.Unauthorized(w, "Invalid or expired token")
				return
			}

			log.Printf("[AUTH] Token valid for user: %s (ID: %d)", claims.Email, claims.UserID) // DEBUG

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext récupère les claims du context
func GetUserFromContext(ctx context.Context) (*auth.JWTClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*auth.JWTClaims)
	return claims, ok
}
