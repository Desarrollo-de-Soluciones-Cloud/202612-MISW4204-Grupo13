package infrastructure

import (
	reportsDomain "backend/internal/reports/domain"
	"backend/internal/shared/database"
	"errors"

	"gorm.io/gorm"
)

type ReportRepository struct{}

var _ reportsDomain.ReportRepository = (*ReportRepository)(nil)

func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

func (r *ReportRepository) Create(report *reportsDomain.Report) error {
	return database.DB.Create(report).Error
}

func (r *ReportRepository) FindByID(id uint) (*reportsDomain.Report, error) {
	var report reportsDomain.Report
	result := database.DB.First(&report, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, reportsDomain.ErrReportNotFound
		}
		return nil, result.Error
	}
	return &report, nil
}

func (r *ReportRepository) FindAll(workspaceID *uint, weekID *uint) ([]reportsDomain.Report, error) {
	var reports []reportsDomain.Report
	query := database.DB.Model(&reportsDomain.Report{}).Order("id desc")
	if workspaceID != nil {
		query = query.Where("workspace_id = ?", *workspaceID)
	}
	if weekID != nil {
		query = query.Where("week_id = ?", *weekID)
	}
	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *ReportRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&reportsDomain.Report{})
}
