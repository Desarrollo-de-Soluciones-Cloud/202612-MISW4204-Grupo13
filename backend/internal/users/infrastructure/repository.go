package infrastructure

import (
	"errors"

	"backend/internal/shared/database"
	"backend/internal/users/domain"

	"gorm.io/gorm"
)

type UserRepository struct{}

var _ domain.UserRepository = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *domain.User) error {
	result := database.DB.Create(user)
	return result.Error
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	result := database.DB.Find(&users)
	return users, result.Error
}

func (r *UserRepository) FindAllByRole(role domain.UserRole) ([]domain.User, error) {
	var users []domain.User
	result := database.DB.Where("global_role = ?", role).Find(&users)
	return users, result.Error
}

func (r *UserRepository) Update(user *domain.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&domain.User{})
}
