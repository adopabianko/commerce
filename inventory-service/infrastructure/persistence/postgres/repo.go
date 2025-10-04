package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "github.com/adopabianko/commerce/inventory-service/internal/domain/inventory"
)

type Repo struct { db *gorm.DB }
func New(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) GetBySKU(ctx context.Context, sku string, forUpdate bool) (*domain.Product, error) {
	var p domain.Product
	q := r.db.WithContext(ctx).Where("sku = ?", sku)
	if forUpdate { q = q.Clauses(clause.Locking{Strength: "UPDATE"}) }
	if err := q.First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
		return nil, err
	}
	return &p, nil
}

func (r *Repo) BulkUpsertProducts(ctx context.Context, ps []domain.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range ps {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "sku"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "stock"}),
			}).Create(&ps[i]).Error; err != nil { return err }
		}
		return nil
	})
}

func (r *Repo) AdjustStock(ctx context.Context, sku string, delta int32) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p domain.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("sku = ?", sku).First(&p).Error; err != nil { return err }
		newStock := p.Stock + delta
		if newStock < 0 { return errors.New("insufficient stock") }
		p.Stock = newStock
		return tx.Save(&p).Error
	})
}
