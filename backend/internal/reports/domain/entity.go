package domain

import "time"

type ReportType string

const ReportTypeWeeklyPDF ReportType = "weekly_pdf"

type Report struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	WorkspaceID  uint       `gorm:"not null;index" json:"workspace_id"`
	WeekID       uint       `gorm:"not null;index" json:"week_id"`
	AssignmentID uint       `gorm:"not null;index" json:"assignment_id"`
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	Type         ReportType `gorm:"size:50;not null" json:"type"`
	Summary      string     `gorm:"type:text;not null" json:"summary"`
	FilePath     string     `gorm:"type:text;not null" json:"file_path"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewWeeklyReport(workspaceID, weekID, assignmentID, userID uint, summary, filePath string) (*Report, error) {
	if workspaceID == 0 {
		return nil, ErrReportWorkspaceIDRequired
	}
	if weekID == 0 {
		return nil, ErrReportWeekIDRequired
	}
	if assignmentID == 0 {
		return nil, ErrReportAssignmentIDRequired
	}
	if userID == 0 {
		return nil, ErrReportUserIDRequired
	}
	if summary == "" {
		return nil, ErrReportSummaryRequired
	}
	if filePath == "" {
		return nil, ErrReportFilePathRequired
	}

	return &Report{
		WorkspaceID:  workspaceID,
		WeekID:       weekID,
		AssignmentID: assignmentID,
		UserID:       userID,
		Type:         ReportTypeWeeklyPDF,
		Summary:      summary,
		FilePath:     filePath,
	}, nil
}
