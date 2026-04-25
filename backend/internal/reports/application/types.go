package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type WorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type WeekReader interface {
	FindByID(id uint) (*weeksDomain.Week, error)
}

type AssignmentReader interface {
	FindAllByWorkspaceID(workspaceID uint) ([]assignmentsDomain.Assignment, error)
}

type TaskReader interface {
	FindAll() ([]tasksDomain.Task, error)
}

type PDFGenerator interface {
	Generate(filePath string, title string, lines []string) error
}

type ReportOutput struct {
	ID           uint   `json:"id"`
	WorkspaceID  uint   `json:"workspace_id"`
	WeekID       uint   `json:"week_id"`
	AssignmentID uint   `json:"assignment_id"`
	UserID       uint   `json:"user_id"`
	Type         string `json:"type"`
	Summary      string `json:"summary"`
	FilePath     string `json:"file_path"`
}

func toReportOutput(report reportsDomain.Report) ReportOutput {
	return ReportOutput{
		ID:           report.ID,
		WorkspaceID:  report.WorkspaceID,
		WeekID:       report.WeekID,
		AssignmentID: report.AssignmentID,
		UserID:       report.UserID,
		Type:         string(report.Type),
		Summary:      report.Summary,
		FilePath:     report.FilePath,
	}
}
