package handlers

import (
	"database/sql"
	"net/http"

	mw "github.com/uit-go/user/internal/httpserver/middleware"
)

func Me(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := mw.UserID(r)
		var email, name, role string
		err := db.QueryRow(`SELECT email, name, role FROM users WHERE id=$1`, uid).
			Scan(&email, &name, &role)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"id": uid, "email": email, "name": name, "role": role,
		})
	}
}
