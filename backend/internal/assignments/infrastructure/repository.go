package infrastructure

import (
	"errors"

	"backend/internal/assignments/domain"
	"backend/internal/shared/database"

	"gorm.io/gorm"
)

type AssignmentRepository struct{}

var _ domain.AssignmentRepository = (*AssignmentRepository)(nil)

func NewAssignmentRepository() *AssignmentRepository {
	return &AssignmentRepository{}
}

func (r *AssignmentRepository) Create(assignment *domain.Assignment) error {
	return database.DB.Create(assignment).Error
}

func (r *AssignmentRepository) FindByID(id uint) (*domain.Assignment, error) {
	var assignment domain.Assignment
	result := database.DB.First(&assignment, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAssignmentNotFound
		}
		return nil, result.Error
	}
	return &assignment, nil
}

func (r *AssignmentRepository) FindAllByUserID(userID uint) ([]domain.Assignment, error) {
	var assignments []domain.Assignment
	result := database.DB.Where("user_id = ?", userID).Find(&assignments)
	return assignments, result.Error
}

func (r *AssignmentRepository) Update(assignment *domain.Assignment) error {
	return database.DB.Save(assignment).Error
}

func (r *AssignmentRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&domain.Assignment{})
}
