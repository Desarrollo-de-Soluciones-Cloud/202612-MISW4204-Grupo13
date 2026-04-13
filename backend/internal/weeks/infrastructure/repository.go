package infrastructure

import (
	"backend/internal/shared/database"
	"backend/internal/weeks/domain"
	"errors"

	"gorm.io/gorm"
)

type WeekRepository struct{}

var _ domain.WeekRepository = (*WeekRepository)(nil)

func NewWeekRepository() *WeekRepository {
	return &WeekRepository{}
}

func (r *WeekRepository) CreateMany(weeks []domain.Week) error {
	return database.DB.Create(&weeks).Error
}

func (r *WeekRepository) FindAllByPeriodID(periodID uint) ([]domain.Week, error) {
	var weeks []domain.Week
	result := database.DB.Where("period_id = ?", periodID).Order("number ASC").Find(&weeks)
	return weeks, result.Error
}

func (r *WeekRepository) ExistsByPeriodID(periodID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&domain.Week{}).Where("period_id = ?", periodID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *WeekRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&domain.Week{})
}

func (r *WeekRepository) FindByID(id uint) (*domain.Week, error) {
	var week domain.Week
	result := database.DB.First(&week, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWeekNotFound
		}
		return nil, result.Error
	}
	return &week, nil
}
