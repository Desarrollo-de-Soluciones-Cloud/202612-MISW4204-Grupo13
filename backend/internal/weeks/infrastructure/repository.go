package infrastructure

import (
	"errors"

	"backend/internal/shared/database"
	"backend/internal/weeks/domain"

	"gorm.io/gorm"
)

type WeekRepository struct{}

var _ domain.WeekRepository = (*WeekRepository)(nil)

func NewWeekRepository() *WeekRepository {
	return &WeekRepository{}
}

func (r *WeekRepository) Create(week *domain.Week) error {
	result := database.DB.Create(week)
	return result.Error
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

func (r *WeekRepository) FindAll() ([]*domain.Week, error) {
	var weeks []*domain.Week
	result := database.DB.Find(&weeks)
	return weeks, result.Error
}

func (r *WeekRepository) FindAllByPeriodID(periodID uint) ([]*domain.Week, error) {
	var weeks []*domain.Week
	result := database.DB.Where("period_id = ?", periodID).Find(&weeks)
	return weeks, result.Error
}

func (r *WeekRepository) Update(week *domain.Week) error {
	return database.DB.Save(week).Error
}

func (r *WeekRepository) Delete(id uint) error {
	result := database.DB.Delete(&domain.Week{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrWeekNotFound
	}
	return nil
}

func (r *WeekRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&domain.Week{})
}
