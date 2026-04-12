package domain

type WeekRepository interface {
	Create(week *Week) error
	FindByID(id uint) (*Week, error)
	FindAll() ([]*Week, error)
	FindAllByPeriodID(periodID uint) ([]*Week, error)
	Update(week *Week) error
	Delete(id uint) error
}
