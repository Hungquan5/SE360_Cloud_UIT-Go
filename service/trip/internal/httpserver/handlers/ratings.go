package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uit-go/trip/internal/httpserver/middleware"
)

// POST /trips/{id}/rate {score, comment, rateeId}
func RateTrip(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tripID := chi.URLParam(r, "id")
		rater := middleware.UserID(r)

		var in struct {
			Score   int    `json:"score"`
			Comment string `json:"comment"`
			RateeID string `json:"rateeId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", 400)
			return
		}

		// ensure trip exists & is completed (simple check)
		var status string
		if err := db.QueryRow(`SELECT status FROM trips WHERE id=$1`, tripID).Scan(&status); err != nil {
			http.Error(w, "trip not found", 404)
			return
		}
		if status != "COMPLETED" {
			http.Error(w, "trip not completed", 409)
			return
		}

		_, err := db.Exec(`
			INSERT INTO ratings (trip_id, rater_id, ratee_id, score, comment)
			VALUES ($1,$2,$3,$4,$5)
		`, tripID, rater, in.RateeID, in.Score, in.Comment)
		if err != nil {
			http.Error(w, "db error", 500)
			return
		}
		w.WriteHeader(201)
	}
}
