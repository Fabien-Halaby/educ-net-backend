package middleware

import (
	"net/http"
	"log"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📨 %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}