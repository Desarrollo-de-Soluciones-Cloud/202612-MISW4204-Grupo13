package domain

import "time"

type Report struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID  uint      `gorm:"not null;uniqueIndex:idx_report_unique" json:"workspace_id"`
	WeekID       uint      `gorm:"not null;uniqueIndex:idx_report_unique" json:"week_id"`
	AssignmentID uint      `gorm:"not null;uniqueIndex:idx_report_unique" json:"assignment_id"`
	UserID       uint      `gorm:"not null;uniqueIndex:idx_report_unique" json:"user_id"`
	FilePath     string    `gorm:"type:text;not null" json:"file_path"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewWeeklyReport(workspaceID, weekID, assignmentID, userID uint, filePath string) (*Report, error) {
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
	if filePath == "" {
		return nil, ErrReportFilePathRequired
	}

	return &Report{
		WorkspaceID:  workspaceID,
		WeekID:       weekID,
		AssignmentID: assignmentID,
		UserID:       userID,
		FilePath:     filePath,
	}, nil
}