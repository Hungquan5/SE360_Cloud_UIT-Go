package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	mw "github.com/uit-go/user/internal/httpserver/middleware"
)

type driverApplyIn struct {
	VehiclePlate string `json:"vehicle_plate"`
	VehicleModel string `json:"vehicle_model"`
}

func DriverApply(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := mw.UserID(r)
		var in driverApplyIn
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		// Upsert simple profile; approval later by admin tool
		_, err := db.Exec(`
			INSERT INTO driver_profiles (user_id, vehicle_plate, vehicle_model, approved, online, created_at, updated_at)
			VALUES ($1,$2,$3,false,false,now(),now())
			ON CONFLICT (user_id) DO UPDATE
			SET vehicle_plate=EXCLUDED.vehicle_plate,
			    vehicle_model=EXCLUDED.vehicle_model,
			    updated_at=now()
		`, uid, in.VehiclePlate, in.VehicleModel)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"status": "submitted"})
	}
}

func DriverMe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := mw.UserID(r)
		var plate, model string
		var approved, online bool
		var updatedAt time.Time

		err := db.QueryRow(`
			SELECT vehicle_plate, vehicle_model, approved, online, updated_at
			FROM driver_profiles WHERE user_id=$1
		`, uid).Scan(&plate, &model, &approved, &online, &updatedAt)
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusOK, map[string]any{"profile": nil})
			return
		}
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"vehicle_plate": plate,
			"vehicle_model": model,
			"approved":      approved,
			"online":        online,
			"updated_at":    updatedAt,
		})
	}
}
