package middleware

import (
	"net/http"
	"strings"
)

var allowedOrigins = []string{
	"http://localhost:3000",
	"http://localhost:8080",
	"https://yourdomain.com",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "https://localhost:") {
		return true
	}

	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}
