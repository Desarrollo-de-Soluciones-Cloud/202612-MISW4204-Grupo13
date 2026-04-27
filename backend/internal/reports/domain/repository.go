package domain

type ReportRepository interface {
	Save(report *Report) error
	FindByID(id uint) (*Report, error)
	FindAll(workspaceID uint, weekID *uint, userID *uint) ([]Report, error)
	AutoMigrate() error
}