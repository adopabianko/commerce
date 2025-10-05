package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBDSN         string
	HTTPAddr      string
	InventoryGRPC string
	ReqTimeout    time.Duration
}

func Load() Config {
	_ = godotenv.Load()
	tm, _ := strconv.Atoi(getenv("REQ_TIMEOUT_MS", "2000"))
	return Config{
		DBDSN:         getenv("DB_DSN", "host=localhost user=postgres password=postgres dbname=order_db port=5432 sslmode=disable TimeZone=UTC"),
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		InventoryGRPC: getenv("INVENTORY_GRPC_ADDR", "localhost:50051"),
		ReqTimeout:    time.Duration(tm) * time.Millisecond,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
