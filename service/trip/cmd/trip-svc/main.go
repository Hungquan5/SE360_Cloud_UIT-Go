package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/uit-go/trip/internal/config"
	"github.com/uit-go/trip/internal/db"
	"github.com/uit-go/trip/internal/events"
	"github.com/uit-go/trip/internal/httpserver"
	"github.com/uit-go/trip/internal/httpserver/ws"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Get()

	// DB
	d, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	// NATS
	bus, err := events.Connect(cfg.NATSURL)
	if err != nil {
		log.Fatalf("nats: %v", err)
	}

	s := &httpserver.Server{
		DB:        d,
		Bus:       bus,
		JWTSecret: cfg.JWTSecret,
		WSHub:     ws.NewHub(),
	}
	mux := s.Router()

	// subscribers
	s.StartSubscribers()

	srv := &http.Server{Addr: cfg.Addr, Handler: mux}
	go func() {
		log.Printf("trip listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	_ = d.Close()
	bus.NATS.Drain()
}
