package trip

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nats-io/nats.go"
	"golang.org/x/net/websocket"
	"time"
)

type Trip struct {
	ID                   string  `json:"id"`
	PassengerID          string  `json:"passenger_id"`
	OriginLat, OriginLng float64 `json:"origin_lat","origin_lng"`
	DestLat, DestLng     float64 `json:"dest_lat","dest_lng"`
	FareEstimateCents    int     `json:"fare_estimate_cents"`
	Status               string  `json:"status"`
}

var (
	db    *sql.DB
	nc    *nats.Conn
	wsHub = NewHub() // broadcast to trip rooms
)

func main() {
	var err error
	db, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
	must(err)
	nc, err = nats.Connect(os.Getenv("NATS_URL"))
	must(err)

	// subscribe driver.location.updated → push to WS
	nc.Subscribe("driver.location.updated", func(m *nats.Msg) {
		wsHub.BroadcastToAll(m.Data) // or to specific trip room if payload has TripID
	})

	r := chi.NewRouter()

	// POST /trips
	r.Post("/trips", func(w http.ResponseWriter, req *http.Request) {
		var in struct{ OriginLat, OriginLng, DestLat, DestLng float64 }
		json.NewDecoder(req.Body).Decode(&in)
		fare := roughFareEstimate(in.OriginLat, in.OriginLng, in.DestLat, in.DestLng)

		var id string
		err := db.QueryRow(`INSERT INTO trips (passenger_id, origin_lat,origin_lng,dest_lat,dest_lng,fare_estimate_cents,status)
      VALUES ($1,$2,$3,$4,$5,$6,'SEARCHING') RETURNING id`,
			currentUser(req), in.OriginLat, in.OriginLng, in.DestLat, in.DestLng, fare).Scan(&id)
		must(err)

		payload, _ := json.Marshal(map[string]any{
			"tripId": id, "originLat": in.OriginLat, "originLng": in.OriginLng,
			"destLat": in.DestLat, "destLng": in.DestLng,
		})
		nc.Publish("trip.requested", payload)

		writeJSON(w, map[string]any{"id": id, "status": "SEARCHING", "fare_estimate_cents": fare})
	})

	// GET /trips/{id}
	r.Get("/trips/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		var t Trip
		err := db.QueryRow(`SELECT id, passenger_id, origin_lat,origin_lng,dest_lat,dest_lng,fare_estimate_cents,status FROM trips WHERE id=$1`, id).
			Scan(&t.ID, &t.PassengerID, &t.OriginLat, &t.OriginLng, &t.DestLat, &t.DestLng, &t.FareEstimateCents, &t.Status)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}
		writeJSON(w, t)
	})

	// POST /trips/{id}/cancel
	r.Post("/trips/{id}/cancel", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		_, err := db.Exec(`UPDATE trips SET status='CANCELLED', updated_at=now() WHERE id=$1 AND status IN ('SEARCHING','ACCEPTED','ONGOING')`, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)
	})

	// WS realtime
	r.Get("/trips/{id}/ws", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		websocket.Handler(func(conn *websocket.Conn) {
			wsHub.Add(id, conn)
			defer wsHub.Remove(id, conn)
			for {
				time.Sleep(time.Hour)
			} // keep open
		}).ServeHTTP(w, req)
	})

	log.Println("trip-svc on :8080")
	http.ListenAndServe(":8080", r)
}
