package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/uit-go/user/internal/security"
)

type registerIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"` // PASSENGER | DRIVER | ADMIN (usually PASSENGER/DRIVER)
}

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in registerIn
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		in.Email = strings.TrimSpace(strings.ToLower(in.Email))
		if in.Email == "" || in.Password == "" || in.Role == "" {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}
		hash, err := security.HashPassword(in.Password)
		if err != nil {
			http.Error(w, "hash error", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`INSERT INTO users (email,password_hash,name,role) VALUES ($1,$2,$3,$4)`,
			in.Email, hash, in.Name, in.Role)
		if err != nil {
			http.Error(w, "email already exists?", http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

type loginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in loginIn
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		in.Email = strings.TrimSpace(strings.ToLower(in.Email))
		var id, role, hash string
		err := db.QueryRow(`SELECT id, role, password_hash FROM users WHERE email=$1`, in.Email).
			Scan(&id, &role, &hash)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if err := security.CheckPassword(hash, in.Password); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		token, err := security.SignHS256(jwtSecret, id, role, 24*time.Hour)
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"access_token": token})
	}
}
