package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/yourorg/driver-svc/internal/cache"
	"github.com/yourorg/driver-svc/internal/httpserver/handlers"
)

func New(redis *cache.Redis, cfg interface{ Port string }) http.Handler {
	log := httplog.NewLogger("driver-svc", httplog.Options{JSON: true})
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(log))

	// health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Get("/ready", handlers.Ready(redis))

	// APIs
	r.Post("/drivers/{id}/location", handlers.UpdateLocation(redis))
	r.Post("/drivers/{id}/status/{state}", handlers.SetStatus(redis))
	r.Get("/drivers/nearby", handlers.Nearby(redis))

	return r
}
