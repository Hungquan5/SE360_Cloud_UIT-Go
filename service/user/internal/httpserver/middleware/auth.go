package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/uit-go/user/internal/security"
)

type ctxKey string

const (
	ctxUserID ctxKey = "uid"
	ctxRole   ctxKey = "role"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			claims, err := security.ParseHS256(secret, token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ctxUserID, claims.Sub)
			ctx = context.WithValue(ctx, ctxRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Helpers to read from context (for handlers)
func UserID(r *http.Request) string {
	if v := r.Context().Value(ctxUserID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func Role(r *http.Request) string {
	if v := r.Context().Value(ctxRole); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
