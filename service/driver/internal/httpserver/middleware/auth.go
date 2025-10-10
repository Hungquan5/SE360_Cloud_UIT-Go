package middleware

import "net/http"

// Placeholder: add JWT verification consistent with user/trip services.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: validate Authorization header, set context values
		next.ServeHTTP(w, r)
	})
}
