package config

import "os"

type Config struct {
	Port      string // APP_PORT (default 8080)
	RedisAddr string // REDIS_ADDR (e.g., "redis:6379" or "redis://:pass@host:6379/0")
}

func getenv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }

func Load() Config {
	return Config{
		Port:      getenv("APP_PORT", "8080"),
		RedisAddr: getenv("REDIS_ADDR", "redis:6379"),
	}
}
