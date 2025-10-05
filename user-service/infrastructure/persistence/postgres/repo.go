package postgres

import (
	"errors"

	"gorm.io/gorm"

	domain "github.com/adopabianko/commerce/user-service/internal/domain/user"
)

type Repo struct{ db *gorm.DB }

func New(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

func (r *Repo) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repo) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
