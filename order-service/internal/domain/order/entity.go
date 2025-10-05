package order

import "time"

type Order struct {
	ID        string `gorm:"primaryKey;size:36"`
	Status    string
	CreatedAt time.Time
	Items     []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type OrderItem struct {
	ID      uint   `gorm:"primaryKey"`
	OrderID string `gorm:"index"`
	SKU     string
	Qty     int32
}

type Repository interface {
	CreateOrder(o *Order) error
	UpdateStatus(id string, status string) error
}
