package config

import (
	"log"
	"os"
)

type Config struct {
	Addr        string // :8080
	DatabaseURL string // postgres://app:password@localhost:5432/appdb?sslmode=disable
	JWTSecret   string // HS256 secret (keep long & random)
	Env         string // dev|prod
}

func Get() Config {
	cfg := Config{
		Addr:        getEnv("ADDR", ":8080"),
		DatabaseURL: mustEnv("DATABASE_URL"),
		JWTSecret:   mustEnv("JWT_SECRET"),
		Env:         getEnv("ENV", "dev"),
	}
	return cfg
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
		log.Fatalf("missing required env: %s", k)
	}
	return v
}
