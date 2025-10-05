package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/adopabianko/commerce/order-service/internal/domain/order"
	inventoryv1 "github.com/adopabianko/commerce/proto/gen/inventory/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type memRepo struct{}

func (m *memRepo) CreateOrder(o *order.Order) error            { return nil }
func (m *memRepo) UpdateStatus(id string, status string) error { return nil }

type fakeInv struct{}

func (f *fakeInv) CheckStock(ctx context.Context, in *inventoryv1.CheckStockRequest, opts ...grpc.CallOption) (*inventoryv1.CheckStockResponse, error) {
	return &inventoryv1.CheckStockResponse{Ok: true}, nil
}
func (f *fakeInv) ReserveStock(ctx context.Context, in *inventoryv1.ReserveStockRequest, opts ...grpc.CallOption) (*inventoryv1.ReserveStockResponse, error) {
	return &inventoryv1.ReserveStockResponse{Ok: true}, nil
}
func (f *fakeInv) ReleaseStock(ctx context.Context, in *inventoryv1.ReleaseStockRequest, opts ...grpc.CallOption) (*inventoryv1.ReleaseStockResponse, error) {
	return &inventoryv1.ReleaseStockResponse{Ok: true}, nil
}

func TestPlaceOrderBasic(t *testing.T) {
	uc := NewPlaceOrder(&memRepo{}, &fakeInv{}, 2*time.Second)
	_, err := uc.Exec(context.Background(), Request{Items: []Item{{SKU: "SKU-1", Qty: 1}}})
	require.NoError(t, err)
}
