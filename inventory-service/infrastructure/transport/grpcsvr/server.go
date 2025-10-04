package grpcsvr

import (
	"context"

	inventoryv1 "github.com/adopabianko/commerce/proto/gen/inventory/v1"
	"github.com/adopabianko/commerce/inventory-service/internal/usecase"
)

type Server struct {
	inventoryv1.UnimplementedInventoryServiceServer
	svc *usecase.Service
}

func NewServer(s *usecase.Service) *Server { return &Server{svc: s} }

func (s *Server) CheckStock(ctx context.Context, req *inventoryv1.CheckStockRequest) (*inventoryv1.CheckStockResponse, error) {
	return s.svc.CheckStock(ctx, req)
}
func (s *Server) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	return s.svc.ReserveStock(ctx, req)
}
func (s *Server) ReleaseStock(ctx context.Context, req *inventoryv1.ReleaseStockRequest) (*inventoryv1.ReleaseStockResponse, error) {
	return s.svc.ReleaseStock(ctx, req)
}
