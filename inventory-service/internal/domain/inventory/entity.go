package inventory

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	SKU   string `gorm:"uniqueIndex;size:64"`
	Name  string
	Stock int32
}
