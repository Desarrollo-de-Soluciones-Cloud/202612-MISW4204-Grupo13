package domain

type ReportRepository interface {
	Create(report *Report) error
	FindByID(id uint) (*Report, error)
	FindAll(workspaceID *uint, weekID *uint) ([]Report, error)
	AutoMigrate() error
}
