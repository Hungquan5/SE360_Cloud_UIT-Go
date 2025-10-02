package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/uit-go/trip/internal/domain"
	"github.com/uit-go/trip/internal/events"
	"github.com/uit-go/trip/internal/httpserver/middleware"
	"github.com/uit-go/trip/internal/httpserver/ws"
)

type TripHandlers struct {
	DB  *sql.DB
	Bus *events.Bus
	WS  *ws.Hub
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // adjust CORS in prod

// POST /trips {originLat,originLng,destLat,destLng}
func (h *TripHandlers) CreateTrip(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserID(r)
	var in struct {
		OriginLat float64 `json:"originLat"`
		OriginLng float64 `json:"originLng"`
		DestLat   float64 `json:"destLat"`
		DestLng   float64 `json:"destLng"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	fare := domain.FareEstimateCents(in.OriginLat, in.OriginLng, in.DestLat, in.DestLng)

	var id string
	err := h.DB.QueryRow(`
		INSERT INTO trips (passenger_id, origin_lat, origin_lng, dest_lat, dest_lng, fare_estimate_cents, status)
		VALUES ($1,$2,$3,$4,$5,$6,'SEARCHING') RETURNING id
	`, uid, in.OriginLat, in.OriginLng, in.DestLat, in.DestLng, fare).Scan(&id)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	// publish trip.requested
	_ = h.Bus.PublishJSON(events.SubTripRequested, domain.TripRequested{
		TripID: id, OriginLat: in.OriginLat, OriginLng: in.OriginLng, DestLat: in.DestLat, DestLng: in.DestLng,
	})

	writeJSON(w, 201, map[string]any{
		"id": id, "status": "SEARCHING", "fare_estimate_cents": fare,
	})
}

// GET /trips/{id}
func (h *TripHandlers) GetTrip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var t domain.Trip
	var driverID sql.NullString
	err := h.DB.QueryRow(`
		SELECT id, passenger_id, driver_id, origin_lat, origin_lng, dest_lat, dest_lng, fare_estimate_cents, status, created_at, updated_at
		FROM trips WHERE id=$1
	`, id).Scan(&t.ID, &t.PassengerID, &driverID, &t.OriginLat, &t.OriginLng, &t.DestLat, &t.DestLng, &t.FareEstimateCents, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	if driverID.Valid {
		v := driverID.String
		t.DriverID = &v
	}
	writeJSON(w, 200, t)
}

// POST /trips/{id}/cancel
func (h *TripHandlers) CancelTrip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.DB.Exec(`
		UPDATE trips SET status='CANCELLED'
		WHERE id=$1 AND status IN ('SEARCHING','ACCEPTED','ONGOING')
	`, id)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	w.WriteHeader(204)
}

// POST /trips/{id}/complete  (driver action)
func (h *TripHandlers) CompleteTrip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.DB.Exec(`UPDATE trips SET status='COMPLETED' WHERE id=$1 AND status IN ('ACCEPTED','ONGOING')`, id)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	// optional publish
	_ = h.Bus.PublishJSON(events.SubTripCompleted, map[string]string{"tripId": id})
	w.WriteHeader(204)
}

// GET /trips/{id}/ws  (passenger subscribes)
func (h *TripHandlers) TripWS(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.WS.Add(tripID, conn)
	defer func() { h.WS.Remove(tripID, conn); _ = conn.Close() }()
	// keep the connection alive; server pushes via h.WS.Broadcast(tripID,...)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		} // read ping/pong/close
	}
}

// Subscriptions from NATS (hook these in server startup)
func (h *TripHandlers) HandleDriverAccepted(raw []byte) {
	var ev domain.DriverAccepted
	if err := json.Unmarshal(raw, &ev); err == nil {
		_, _ = h.DB.Exec(`UPDATE trips SET driver_id=$1, status='ACCEPTED' WHERE id=$2 AND status='SEARCHING'`, ev.DriverID, ev.TripID)
		// You could also broadcast a status update to passengers here.
	}
}

func (h *TripHandlers) HandleDriverLocation(raw []byte) {
	// raw message is broadcast to the trip room if TripID present
	h.WS.BroadcastAll(raw) // or decode and route by tripId:
	// var loc domain.DriverLocation; _ = json.Unmarshal(raw,&loc)
	// if loc.TripID != "" { h.WS.Broadcast(loc.TripID, raw) }
}

// (optional) mark trip as ONGOING when driver reaches pickup
func (h *TripHandlers) SetOngoing(tripID string) {
	_, _ = h.DB.Exec(`UPDATE trips SET status='ONGOING', updated_at=now() WHERE id=$1 AND status='ACCEPTED'`, tripID)
}

// helper
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
