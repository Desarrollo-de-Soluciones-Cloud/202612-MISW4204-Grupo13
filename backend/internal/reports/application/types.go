package application

import (
	"context"
	"io"
	"time"

	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type WorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type WeekReader interface {
	FindByID(id uint) (*weeksDomain.Week, error)
}

type UserReader interface {
	FindByID(id uint) (*usersDomain.User, error)
}

type AssignmentReader interface {
	FindAllByWorkspaceID(workspaceID uint) ([]assignmentsDomain.Assignment, error)
}

type TaskReader interface {
	FindAllByWorkspaceAndWeek(workspaceID uint, weekID uint, weekInitialDate string) ([]tasksDomain.Task, error)
}

type PDFGenerator interface {
	Generate(filePath string, title string, lines []string) error
}

type ReportFileStorage interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, contentType string) error
}

type AIReportGenerator interface {
	GenerateWeeklyReport(input AIWeeklyReportInput) (string, error)
}

type AIWeeklyReportInput struct {
	WorkspaceID   uint
	WorkspaceName string
	WeekID        uint
	WeekNumber    int
	InitialDate   string
	FinalDate     string
	AssignmentID  uint
	UserID        uint
	UserName      string
	Role          string
	TotalHours    int
	Tasks         []AIWeeklyReportTask
}

type AIWeeklyReportTask struct {
	Title        string
	Description  string
	Status       string
	SpentHours   int
	Observations string
	Late         bool
}

type ReportOutput struct {
	ID           uint      `json:"id"`
	WorkspaceID  uint      `json:"workspace_id"`
	WeekID       uint      `json:"week_id"`
	AssignmentID uint      `json:"assignment_id"`
	UserID       uint      `json:"user_id"`
	FilePath     string    `json:"file_path"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func toReportOutput(report reportsDomain.Report) ReportOutput {
	return ReportOutput{
		ID:           report.ID,
		WorkspaceID:  report.WorkspaceID,
		WeekID:       report.WeekID,
		AssignmentID: report.AssignmentID,
		UserID:       report.UserID,
		FilePath:     report.FilePath,
		CreatedAt:    report.CreatedAt,
		UpdatedAt:    report.UpdatedAt,
	}
}