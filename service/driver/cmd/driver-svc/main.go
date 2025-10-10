package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourorg/driver-svc/internal/cache"
	"github.com/yourorg/driver-svc/internal/config"
	"github.com/yourorg/driver-svc/internal/httpserver"
)

func main() {
	cfg := config.Load()

	redisStore, err := cache.NewRedis(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis connect failed: %v", err)
	}

	srv := httpserver.New(redisStore, cfg)
	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("driver-svc listening :%s (redis=%s)", cfg.Port, cfg.RedisAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
}
