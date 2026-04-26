package delivery

import "time"

type GenerateWeeklyReportsRequest struct {
	WorkspaceID uint `json:"workspace_id" binding:"required"`
	WeekID      uint `json:"week_id" binding:"required"`
}

type ReportResponse struct {
	ID           uint      `json:"id"`
	WorkspaceID  uint      `json:"workspace_id"`
	WeekID       uint      `json:"week_id"`
	AssignmentID uint      `json:"assignment_id"`
	UserID       uint      `json:"user_id"`
	FilePath     string    `json:"file_path"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GenerateWeeklyReportsResponse struct {
	Reports        []ReportResponse `json:"reports"`
	GeneratedCount int              `json:"generated_count"`
}

type ListReportsResponse struct {
	Reports []ReportResponse `json:"reports"`
}