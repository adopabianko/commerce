package postgres

import (
	"github.com/adopabianko/commerce/order-service/internal/domain/order"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func New(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) CreateOrder(o *order.Order) error { return r.db.Create(o).Error }
func (r *Repo) UpdateStatus(id string, status string) error {
	return r.db.Model(&order.Order{}).Where("id = ?", id).Update("status", status).Error
}
