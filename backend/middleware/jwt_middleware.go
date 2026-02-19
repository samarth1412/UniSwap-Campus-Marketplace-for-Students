package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type tokenParser interface {
	ParseToken(tokenString string) (int64, error)
}

type contextKey string

const userIDContextKey contextKey = "user_id"

func Auth(parser tokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("auth_middleware: request started method=%s path=%s", r.Method, r.URL.Path)

			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				log.Printf("auth_middleware: missing authorization header method=%s path=%s", r.Method, r.URL.Path)
				writeUnauthorized(w, "authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || strings.TrimSpace(parts[1]) == "" {
				log.Printf("auth_middleware: invalid authorization format method=%s path=%s", r.Method, r.URL.Path)
				writeUnauthorized(w, "invalid authorization header format")
				return
			}

			userID, err := parser.ParseToken(strings.TrimSpace(parts[1]))
			if err != nil {
				log.Printf("auth_middleware: token parse failed method=%s path=%s err=%v", r.Method, r.URL.Path, err)
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			log.Printf("auth_middleware: token validated user_id=%d method=%s path=%s", userID, r.Method, r.URL.Path)
			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	return userID, ok
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"success":false,"error":"` + message + `"}`))
}
