package middleware

import (
	"net/http"

	"backend-porto/internal/shared/response"
)

// APIKey returns middleware that requires a matching X-API-Key header.
func APIKey(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublic(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			if r.Header.Get("X-API-Key") != expected {
				response.Error(w, http.StatusUnauthorized, "invalid api key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// isPublic returns true for routes reachable via browser navigation
// (health probe and GitHub OAuth redirect targets that cannot send headers).
func isPublic(path string) bool {
	switch path {
	case "/health", "/api/auth/github/login", "/api/auth/github/callback":
		return true
	}
	return false
}
