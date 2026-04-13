package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/shar1mo/wishlist-api/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type JWTParser interface {
	Parse(tokenString string) (*auth.Claims, error)
}

func Auth(jwtManager JWTParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeUnauthorized(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
				writeUnauthorized(w, "invalid authorization header")
				return
			}

			claims, err := jwtManager.Parse(parts[1])
			if err != nil {
				writeUnauthorized(w, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"` + message + `"}`))
}