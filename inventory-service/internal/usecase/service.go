package usecase

import (
	"context"
	"fmt"

	domain "github.com/adopabianko/commerce/inventory-service/internal/domain/inventory"
	inventoryv1 "github.com/adopabianko/commerce/proto/gen/inventory/v1"
)

type Service struct {
	repo domain.Repository
}

func New(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) CheckStock(ctx context.Context, req *inventoryv1.CheckStockRequest) (*inventoryv1.CheckStockResponse, error) {
	short := map[string]int32{}
	for _, it := range req.Items {
		p, _ := s.repo.GetBySKU(ctx, it.Sku, false)
		if p == nil || p.Stock < it.Qty {
			missing := it.Qty
			if p != nil {
				missing = it.Qty - p.Stock
			}
			short[it.Sku] = missing
		}
	}
	ok := len(short) == 0
	msg := "ok"
	if !ok {
		msg = "insufficient stock"
	}
	return &inventoryv1.CheckStockResponse{Ok: ok, Message: msg, Shortages: short}, nil
}

func (s *Service) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	for _, it := range req.Items {
		if err := s.repo.AdjustStock(ctx, it.Sku, -it.Qty); err != nil {
			return &inventoryv1.ReserveStockResponse{Ok: false, Message: fmt.Sprintf("reserve failed: %v", err)}, nil
		}
	}
	return &inventoryv1.ReserveStockResponse{Ok: true, Message: "reserved"}, nil
}

func (s *Service) ReleaseStock(ctx context.Context, req *inventoryv1.ReleaseStockRequest) (*inventoryv1.ReleaseStockResponse, error) {
	for _, it := range req.Items {
		_ = s.repo.AdjustStock(ctx, it.Sku, it.Qty)
	}
	return &inventoryv1.ReleaseStockResponse{Ok: true, Message: "released"}, nil
}
