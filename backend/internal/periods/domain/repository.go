package domain

type PeriodRepository interface {
	Create(period *Period) error
	FindByID(id uint) (*Period, error)
	FindByName(name string) (*Period, error)
	FindAll() ([]Period, error)
	FindAllByState(state PeriodState) ([]Period, error)
	Update(period *Period) error
}
