package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/adopabianko/commerce/order-service/infrastructure/auth"
	pgrepo "github.com/adopabianko/commerce/order-service/infrastructure/persistence/postgres"
	grpccli "github.com/adopabianko/commerce/order-service/infrastructure/transport/grpcclient"
	httpapi "github.com/adopabianko/commerce/order-service/infrastructure/transport/httpapi"
	"github.com/adopabianko/commerce/order-service/internal/config"
	"github.com/adopabianko/commerce/order-service/internal/domain/order"
	"github.com/adopabianko/commerce/order-service/internal/usecase"
)

func main() {
	cfg := config.Load()

	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("DB not ready, retry in 2s... (%v)", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	if err := db.AutoMigrate(&order.Order{}, &order.OrderItem{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	inv, err := grpccli.New(cfg.InventoryGRPC)
	if err != nil {
		log.Fatalf("grpc dial inventory: %v", err)
	}

	repo := pgrepo.New(db)
	uc := usecase.NewPlaceOrder(repo, inv, cfg.ReqTimeout)

	r := gin.Default()
	h := httpapi.New(uc)

	authClient, err := auth.NewGRPCAuthClient()
	if err != nil {
		log.Fatalf("failed to connect user-service grpc: %v", err)
	}
	defer authClient.Close()

	h.Routes(r, authClient)

	log.Printf("order-service HTTP %s", cfg.HTTPAddr)
	_ = r.Run(cfg.HTTPAddr)
}
