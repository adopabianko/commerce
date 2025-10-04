package main

import (
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	inventoryv1 "github.com/adopabianko/commerce/proto/gen/inventory/v1"
	"github.com/adopabianko/commerce/inventory-service/internal/config"
	domain "github.com/adopabianko/commerce/inventory-service/internal/domain/inventory"
	"github.com/adopabianko/commerce/inventory-service/internal/usecase"
	pgrepo "github.com/adopabianko/commerce/inventory-service/infrastructure/persistence/postgres"
	grpcsvr "github.com/adopabianko/commerce/inventory-service/infrastructure/transport/grpcsvr"
	httpadmin "github.com/adopabianko/commerce/inventory-service/infrastructure/transport/httpadmin"
)

func main() {
	cfg := config.Load()

	var db *gorm.DB
	var err error
	for i:=0; i<cfg.DBMaxRetry; i++ {
		db, err = gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
		if err == nil { break }
		log.Printf("DB not ready, retrying in %v... (%v)", cfg.DBRetryWait, err)
		time.Sleep(cfg.DBRetryWait)
	}
	if err != nil { log.Fatalf("db open failed after retries: %v", err) }
	if err := db.AutoMigrate(&domain.Product{}); err != nil { log.Fatalf("migrate: %v", err) }

	repo := pgrepo.New(db)
	svc := usecase.New(repo)

	lis, err := net.Listen("tcp", cfg.GRPCAddr); if err != nil { log.Fatalf("listen: %v", err) }
	grpcServer := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(grpcServer, grpcsvr.NewServer(svc))
	reflection.Register(grpcServer)
	go func(){ log.Printf("gRPC %s", cfg.GRPCAddr); if err := grpcServer.Serve(lis); err != nil { log.Fatalf("grpc serve: %v", err)} }()

	r := gin.Default()
	httpadmin.Router(r, repo)
	go func(){ log.Printf("admin http %s", cfg.AdminHTTP); _ = r.Run(cfg.AdminHTTP) }()

	select {}
}
