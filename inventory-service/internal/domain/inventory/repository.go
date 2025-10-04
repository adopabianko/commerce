package inventory

import "context"

type Repository interface {
	GetBySKU(ctx context.Context, sku string, forUpdate bool) (*Product, error)
	BulkUpsertProducts(ctx context.Context, ps []Product) error
	AdjustStock(ctx context.Context, sku string, delta int32) error
}
