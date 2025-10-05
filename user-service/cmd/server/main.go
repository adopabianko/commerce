package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	userv1 "github.com/adopabianko/commerce/proto/gen/user/v1"
	pgrepo "github.com/adopabianko/commerce/user-service/infrastructure/persistence/postgres"
	grpcsvr "github.com/adopabianko/commerce/user-service/infrastructure/transport/grpcsvr"
	httpsvr "github.com/adopabianko/commerce/user-service/infrastructure/transport/httpsvr"
	"github.com/adopabianko/commerce/user-service/internal/config"
	domain "github.com/adopabianko/commerce/user-service/internal/domain/user"
	"github.com/adopabianko/commerce/user-service/internal/usecase"
)

func main() {
	cfg := config.Load()

	var db *gorm.DB
	var err error
	for range cfg.DBMaxRetry {
		db, err = gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("DB not ready, retrying in %v... (%v)", cfg.DBRetryWait, err)
		time.Sleep(cfg.DBRetryWait)
	}
	if err != nil {
		log.Fatalf("db open failed after retries: %v", err)
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := pgrepo.New(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}
	svc := usecase.New(repo, jwtSecret)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, grpcsvr.NewServer(svc))
	reflection.Register(grpcServer)
	go func() {
		log.Printf("gRPC %s", cfg.GRPCAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	r := gin.Default()
	httpsvr.Router(r, svc, repo, jwtSecret)
	go func() { log.Printf("admin http %s", cfg.HTTPAddr); _ = r.Run(cfg.HTTPAddr) }()

	select {}
}
