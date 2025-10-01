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
	"github.com/uit-go/user/internal/config"
	"github.com/uit-go/user/internal/db"
	"github.com/uit-go/user/internal/httpserver"
)

func main() {
	_ = godotenv.Load() // load .env in dev

	cfg := config.Get()
	d, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}

	s := &httpserver.Server{DB: d, JWTSecret: cfg.JWTSecret}
	srv := &http.Server{Addr: cfg.Addr, Handler: s.Router()}

	// graceful shutdown
	go func() {
		log.Printf("user listening on %s (env=%s)", cfg.Addr, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	_ = d.Close()
}
