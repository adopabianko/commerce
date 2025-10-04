package config

import (
	"os"
	"time"
	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN       string
	GRPCAddr    string
	AdminHTTP   string
	DBMaxRetry  int
	DBRetryWait time.Duration
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		DBDSN:       getenv("DB_DSN", "host=localhost user=postgres password=postgres dbname=inventory_db port=5432 sslmode=disable TimeZone=UTC"),
		GRPCAddr:    getenv("GRPC_ADDR", ":50051"),
		AdminHTTP:   getenv("ADMIN_HTTP_ADDR", ":9090"),
		DBMaxRetry:  10,
		DBRetryWait: 2 * time.Second,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}
