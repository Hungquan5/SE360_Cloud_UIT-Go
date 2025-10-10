package httpserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	h "github.com/uit-go/user/internal/httpserver/handlers"
	mw "github.com/uit-go/user/internal/httpserver/middleware"
)

type Server struct {
	DB        *sql.DB
	JWTSecret string
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// healthz
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/", func(root chi.Router) {
		// public routes (no auth)
		root.Post("/users", h.Register(s.DB))
		root.Post("/sessions", h.Login(s.DB, s.JWTSecret))

		// protected routes (with auth)
		root.Group(func(pr chi.Router) {
			pr.Use(mw.Auth(s.JWTSecret))
			pr.Get("/users/me", h.Me(s.DB))
			pr.Post("/drivers/apply", h.DriverApply(s.DB))
			pr.Get("/drivers/me", h.DriverMe(s.DB))
		})
	})

	return r
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func logErr(err error) {
	if err != nil {
		log.Println("error:", err)
	}
}
