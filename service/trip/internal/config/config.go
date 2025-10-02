package config

import (
	"log"
	"os"
)

type Config struct {
	Addr        string // :8080
	DatabaseURL string
	JWTSecret   string // HS256 for MVP
	NATSURL     string // nats://localhost:4222
	Env         string // dev|prod
}

func Get() Config {
	return Config{
		Addr:        getEnv("ADDR", ":8080"),
		DatabaseURL: mustEnv("DATABASE_URL"),
		JWTSecret:   mustEnv("JWT_SECRET"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		Env:         getEnv("ENV", "dev"),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}
