package application_test

import (
	applicationpkg "backend/internal/reports/application"
	reportsDomain "backend/internal/reports/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
)

func TestListReportsRejectsMissingWorkspaceFilter(t *testing.T) {
	useCase := applicationpkg.NewListReports(
		newMockReportRepository(),
		&mockWorkspaceReader{},
		&mockWeekReader{},
		&mockUserReader{},
	)

	_, err := useCase.Execute(applicationpkg.ListReportsInput{})
	if !errors.Is(err, reportsDomain.ErrReportWorkspaceFilterRequired) {
		t.Fatalf("expected ErrReportWorkspaceFilterRequired, got %v", err)
	}
}

func TestListReportsRejectsUnknownWeekFilter(t *testing.T) {
	weekID := uint(8)
	useCase := applicationpkg.NewListReports(
		newMockReportRepository(),
		&mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1}},
		&mockWeekReader{err: errors.New("missing")},
		&mockUserReader{},
	)

	_, err := useCase.Execute(applicationpkg.ListReportsInput{WorkspaceID: 1, WeekID: &weekID})
	if !errors.Is(err, reportsDomain.ErrReportWeekNotFound) {
		t.Fatalf("expected ErrReportWeekNotFound, got %v", err)
	}
}

func TestListReportsSuccess(t *testing.T) {
	repo := newMockReportRepository()
	repo.reports[1] = &reportsDomain.Report{ID: 1, WorkspaceID: 1, WeekID: 2, AssignmentID: 3, UserID: 4, FilePath: "reports/a.pdf"}
	repo.reports[2] = &reportsDomain.Report{ID: 2, WorkspaceID: 1, WeekID: 2, AssignmentID: 5, UserID: 6, FilePath: "reports/b.pdf"}
	weekID := uint(2)
	userID := uint(4)

	useCase := applicationpkg.NewListReports(
		repo,
		&mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1}},
		&mockWeekReader{week: &weeksDomain.Week{ID: 2}},
		&mockUserReader{user: &usersDomain.User{ID: 4}},
	)

	output, err := useCase.Execute(applicationpkg.ListReportsInput{WorkspaceID: 1, WeekID: &weekID, UserID: &userID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(output.Reports))
	}
	if output.Reports[0].UserID != 4 {
		t.Fatalf("expected filtered user 4, got %d", output.Reports[0].UserID)
	}
}
