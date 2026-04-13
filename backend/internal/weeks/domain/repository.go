package domain

type WeekRepository interface {
	CreateMany(weeks []Week) error
	FindAllByPeriodID(periodID uint) ([]Week, error)
	ExistsByPeriodID(periodID uint) (bool, error)
}
