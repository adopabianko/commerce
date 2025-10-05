package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN       string
	GRPCAddr    string
	HTTPAddr    string
	DBMaxRetry  int
	DBRetryWait time.Duration
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		DBDSN:       getenv("DB_DSN", "host=localhost user=postgres password=postgres dbname=user_db port=5432 sslmode=disable TimeZone=UTC"),
		GRPCAddr:    getenv("GRPC_ADDR", ":50052"),
		HTTPAddr:    getenv("HTTP_ADDR", ":7070"),
		DBMaxRetry:  10,
		DBRetryWait: 2 * time.Second,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
