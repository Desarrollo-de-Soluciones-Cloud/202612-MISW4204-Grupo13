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

func (r *ReportRepository) Save(report *reportsDomain.Report) error {
	var existingReport reportsDomain.Report

	result := database.DB.
		Where(
			"workspace_id = ? AND week_id = ? AND assignment_id = ? AND user_id = ?",
			report.WorkspaceID,
			report.WeekID,
			report.AssignmentID,
			report.UserID,
		).
		First(&existingReport)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return database.DB.Create(report).Error
		}

		return result.Error
	}

	existingReport.FilePath = report.FilePath

	if err := database.DB.Save(&existingReport).Error; err != nil {
		return err
	}

	*report = existingReport

	return nil
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

func (r *ReportRepository) FindAll(workspaceID uint, weekID *uint, userID *uint) ([]reportsDomain.Report, error) {
	var reports []reportsDomain.Report

	query := database.DB.
		Model(&reportsDomain.Report{}).
		Where("workspace_id = ?", workspaceID).
		Order("id desc")

	if weekID != nil {
		query = query.Where("week_id = ?", *weekID)
	}

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *ReportRepository) AutoMigrate() error {
	if err := database.DB.AutoMigrate(&reportsDomain.Report{}); err != nil {
		return err
	}

	if err := database.DB.Exec(`ALTER TABLE reports DROP COLUMN IF EXISTS "type"`).Error; err != nil {
		return err
	}

	if err := database.DB.Exec(`ALTER TABLE reports DROP COLUMN IF EXISTS summary`).Error; err != nil {
		return err
	}

	return nil
}