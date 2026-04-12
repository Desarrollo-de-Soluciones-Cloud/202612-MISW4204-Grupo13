package infrastructure

import (
	"errors"

	"backend/internal/periods/domain"
	"backend/internal/shared/database"

	"gorm.io/gorm"
)

type PeriodRepository struct{}

var _ domain.PeriodRepository = (*PeriodRepository)(nil)

func NewPeriodRepository() *PeriodRepository {
	return &PeriodRepository{}
}

func (r *PeriodRepository) Create(period *domain.Period) error {
	result := database.DB.Create(period)
	return result.Error
}

func (r *PeriodRepository) FindByID(id uint) (*domain.Period, error) {
	var period domain.Period
	result := database.DB.First(&period, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPeriodNotFound
		}
		return nil, result.Error
	}
	return &period, nil
}

func (r *PeriodRepository) FindByName(name string) (*domain.Period, error) {
	var period domain.Period
	result := database.DB.Where("name = ?", name).First(&period)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPeriodNotFound
		}
		return nil, result.Error
	}
	return &period, nil
}

func (r *PeriodRepository) FindAll() ([]domain.Period, error) {
	var periods []domain.Period
	result := database.DB.Find(&periods)
	return periods, result.Error
}

func (r *PeriodRepository) FindAllByState(state domain.PeriodState) ([]domain.Period, error) {
	var periods []domain.Period
	result := database.DB.Where("period_state = ?", state).Find(&periods)
	return periods, result.Error
}

func (r *PeriodRepository) Update(period *domain.Period) error {
	return database.DB.Save(period).Error
}

func (r *PeriodRepository) Delete(id uint) error {
	result := database.DB.Delete(&domain.Period{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPeriodNotFound
	}
	return nil
}

func (r *PeriodRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&domain.Period{})
}
