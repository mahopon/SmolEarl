package main

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs all HTTP requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.RemoteAddr, r.Method, wrapper.statusCode, duration)

	})
}

// CORSMiddleware adds Cross-Origin Resource Sharing headers to each response.
// It also handles preflight OPTIONS requests by returning the appropriate
// headers without invoking the next handler.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set common CORS headers. Adjust the allowed origins, methods, and
		// headers as needed for the application.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// If this is a preflight request, respond with 200 OK and do not
		// forward the request to the next handler.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, continue processing the request.
		next.ServeHTTP(w, r)
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterWrapper) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
