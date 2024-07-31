package middleware

import (
	"github.com/mbswe/selma"
	"net/http"
)

// AuthMiddleware checks for a valid authentication token
func AuthMiddleware(app *selma.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Example auth check
			if r.Header.Get("Authorization") == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
