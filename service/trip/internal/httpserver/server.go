package httpserver

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/uit-go/trip/internal/events"
	"github.com/uit-go/trip/internal/httpserver/handlers"
	mw "github.com/uit-go/trip/internal/httpserver/middleware"
	"github.com/uit-go/trip/internal/httpserver/ws"
)

type Server struct {
	DB        *sql.DB
	Bus       *events.Bus
	JWTSecret string
	WSHub     *ws.Hub
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// healthz
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// protected APIs
	r.Group(func(pr chi.Router) {
		pr.Use(mw.Auth(s.JWTSecret))

		th := &handlers.TripHandlers{DB: s.DB, Bus: s.Bus, WS: s.WSHub}

		pr.Post("/trips", th.CreateTrip)
		pr.Get("/trips/{id}", th.GetTrip)
		pr.Post("/trips/{id}/cancel", th.CancelTrip)
		pr.Post("/trips/{id}/complete", th.CompleteTrip)
		pr.Get("/trips/{id}/ws", th.TripWS)

		pr.Post("/trips/{id}/rate", handlers.RateTrip(s.DB))

		// NATS subscriptions wired below in StartSubscribers()
	})

	return r
}

func (s *Server) StartSubscribers() {
	th := &handlers.TripHandlers{DB: s.DB, Bus: s.Bus, WS: s.WSHub}
	// driver.accepted → update trip to ACCEPTED
	sub1, err := s.Bus.Subscribe(events.SubDriverAccepted, th.HandleDriverAccepted)
	events.LogErr(err)
	// driver.location.updated → broadcast WS
	sub2, err := s.Bus.Subscribe(events.SubDriverLocation, th.HandleDriverLocation)
	events.LogErr(err)
	log.Printf("NATS subscribers ready: %v %v", sub1 != nil, sub2 != nil)
}
