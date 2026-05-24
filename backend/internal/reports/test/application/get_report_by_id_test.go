package application_test

import (
	applicationpkg "backend/internal/reports/application"
	reportsDomain "backend/internal/reports/domain"
	"errors"
	"testing"
)

func TestGetReportByIDSuccess(t *testing.T) {
	repo := newMockReportRepository()
	repo.reports[1] = &reportsDomain.Report{ID: 1, WorkspaceID: 10, WeekID: 20, AssignmentID: 30, UserID: 40, FilePath: "reports/file.pdf"}

	useCase := applicationpkg.NewGetReportByID(repo)
	output, err := useCase.Execute(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.ID != 1 || output.FilePath != "reports/file.pdf" {
		t.Fatalf("unexpected output: %+v", output)
	}
}

func TestGetReportByIDNotFound(t *testing.T) {
	useCase := applicationpkg.NewGetReportByID(newMockReportRepository())

	_, err := useCase.Execute(99)
	if !errors.Is(err, reportsDomain.ErrReportNotFound) {
		t.Fatalf("expected ErrReportNotFound, got %v", err)
	}
}
