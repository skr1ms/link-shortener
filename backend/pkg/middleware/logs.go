package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)
		duration := time.Since(startTime)
		log.Printf("\nRequest: %d %s %s %s %s", wrapper.statusCode, r.Method, r.URL.Path, r.RemoteAddr, duration)
	})
}
