package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"backend-porto/internal/shared/constants"
	"backend-porto/internal/shared/response"

	"github.com/golang-jwt/jwt/v5"
)

// Auth creates middleware that validates JWT Bearer tokens.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractToken(r)
			if tokenString == "" {
				response.Error(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Validate expiration
			exp, err := claims.GetExpirationTime()
			if err != nil || exp.Before(time.Now()) {
				response.Error(w, http.StatusUnauthorized, "token has expired")
				return
			}

			// Extract user info from claims
			userID, _ := claims.GetSubject()
			email, _ := claims.GetIssuer()

			ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
			ctx = context.WithValue(ctx, constants.UserEmailKey, email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth creates middleware that extracts JWT if present but doesn't block.
func OptionalAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractToken(r)
			if tokenString == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err == nil && token.Valid {
				userID, _ := claims.GetSubject()
				email, _ := claims.GetIssuer()

				ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
				ctx = context.WithValue(ctx, constants.UserEmailKey, email)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Also check cookie
		cookie, err := r.Cookie("porto_token")
		if err == nil && cookie.Value != "" {
			return cookie.Value
		}
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}

	return parts[1]
}
