package middleware

import (
	"github.com/mbswe/selma"
	"net/http"
)

// responseWriter is a custom http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	//rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs each request including the status code
func LoggingMiddleware(app *selma.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(rw, r)
			app.MiddlewareLogger.Printf("Request: %s %s, Status: %d", r.Method, r.URL.Path, rw.statusCode)

			if app.Config.Mode == "development" {
				app.DebugLogger.Printf("Request: %s %s, Headers: %s, Status: %d", r.Method, r.URL.Path, r.Header, rw.statusCode)
			}
		})
	}
}
