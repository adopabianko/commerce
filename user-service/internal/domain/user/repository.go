package inventory

type Repository interface {
	Create(u *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}
