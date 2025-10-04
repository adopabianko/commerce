package inventory

import (
	"context"
	"google.golang.org/grpc"
	inventoryv1 "github.com/adopabianko/commerce/proto/gen/inventory/v1"
)

type Client interface {
	CheckStock(ctx context.Context, in *inventoryv1.CheckStockRequest, opts ...grpc.CallOption) (*inventoryv1.CheckStockResponse, error)
	ReserveStock(ctx context.Context, in *inventoryv1.ReserveStockRequest, opts ...grpc.CallOption) (*inventoryv1.ReserveStockResponse, error)
	ReleaseStock(ctx context.Context, in *inventoryv1.ReleaseStockRequest, opts ...grpc.CallOption) (*inventoryv1.ReleaseStockResponse, error)
}
