package main

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	infra_prom "github.com/mahopon/SmolEarl/infra/prometheus"
)

// LoggingMiddleware logs all HTTP requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		logAttrs := []any{
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapper.statusCode),
			slog.Duration("duration", duration),
		}

		switch {
		case wrapper.statusCode >= 500:
			slog.Error("request completed with server error", logAttrs...)
		case wrapper.statusCode >= 400:
			slog.Warn("request completed with client error", logAttrs...)
		case wrapper.statusCode >= 300:
			slog.Info("request completed with redirect", logAttrs...)
		default:
			slog.Debug("request completed successfully", logAttrs...)
		}

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

func StripTrailingSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			r.URL.Path = path[:len(path)-1]
		}
		next.ServeHTTP(w, r)
	})
}

func PrometheusHTTPMiddleware(httpMetrics *infra_prom.HTTPMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			rw := &responseWriterWrapper{w, http.StatusOK}
			next.ServeHTTP(rw, r)
			if !strings.Contains(path, "metrics") {
				httpMetrics.TotalRequests.WithLabelValues(
					strconv.Itoa(rw.statusCode), r.Method,
				).Inc()
			}
		})
	}
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
