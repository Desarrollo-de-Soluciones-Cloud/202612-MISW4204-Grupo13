package domain

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]User, error)
}