package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yourorg/driver-svc/internal/cache"
	"github.com/yourorg/driver-svc/internal/domain"
)

func Ready(redis *cache.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := redis.Ready(ctx); err != nil {
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func SetStatus(redis *cache.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		state := domain.DriverStatus(chi.URLParam(r, "state"))
		if state != domain.StatusOnline && state != domain.StatusOffline {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		if err := redis.SetDriverStatus(r.Context(), id, state); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateLocation(redis *cache.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req domain.UpdateLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", 400)
			return
		}
		if req.Status == "" {
			req.Status = domain.StatusOnline
		}
		if err := redis.UpdateLocation(r.Context(), id, req.Lat, req.Lon); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = redis.SetDriverStatus(r.Context(), id, req.Status)
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"driver_id": id, "lat": req.Lat, "lon": req.Lon,
			"status": req.Status, "at": time.Now().UTC().Format(time.RFC3339Nano),
		})
	}
}

func Nearby(redis *cache.Redis) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		lat, err := strconv.ParseFloat(q.Get("lat"), 64)
		if err != nil {
			http.Error(w, "lat required", 400)
			return
		}
		lon, err := strconv.ParseFloat(q.Get("lon"), 64)
		if err != nil {
			http.Error(w, "lon required", 400)
			return
		}

		radius := 1500.0
		if s := q.Get("radius_m"); s != "" {
			if radius, err = strconv.ParseFloat(s, 64); err != nil || radius <= 0 {
				http.Error(w, "invalid radius_m", 400)
				return
			}
		}
		limit := 50
		if s := q.Get("limit"); s != "" {
			if i, e := strconv.Atoi(s); e == nil && i > 0 {
				limit = i
			}
		}

		res, err := redis.Nearby(r.Context(), lat, lon, radius, limit)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}
}
