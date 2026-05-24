package application_test

import (
	applicationpkg "backend/internal/reports/application"
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
	"time"
)

const (
	testWeekInitialDate = "2026-04-07"
	testWeekFinalDate   = "2026-04-13"
)

func TestGenerateWeeklyReportsRejectsMissingWorkspaceID(t *testing.T) {
	useCase := applicationpkg.NewGenerateWeeklyReports(
		applicationpkg.GenerateWeeklyReportsDependencies{
			ReportRepo:        newMockReportRepository(),
			WorkspaceReader:   &mockWorkspaceReader{},
			WeekReader:        &mockWeekReader{},
			AssignmentReader:  &mockAssignmentReader{},
			TaskReader:        &mockTaskReader{},
			UserReader:        &mockUserReader{},
			PDFGenerator:      &mockPDFGenerator{},
			AIReportGenerator: &mockAIReportGenerator{},
		},
		nil,
	)

	_, err := useCase.Execute(applicationpkg.GenerateWeeklyReportsInput{WeekID: 1})
	if !errors.Is(err, reportsDomain.ErrReportWorkspaceIDRequired) {
		t.Fatalf("expected ErrReportWorkspaceIDRequired, got %v", err)
	}
}

func TestGenerateWeeklyReportsRejectsWhenNoAssignmentsAreReportable(t *testing.T) {
	useCase := applicationpkg.NewGenerateWeeklyReports(
		applicationpkg.GenerateWeeklyReportsDependencies{
			ReportRepo:      newMockReportRepository(),
			WorkspaceReader: &mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "WS"}},
			WeekReader: &mockWeekReader{week: &weeksDomain.Week{
				ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate,
			}},
			AssignmentReader:  &mockAssignmentReader{assignments: []assignmentsDomain.Assignment{{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.AssignmentRole("professor"), WeeklyHours: 4}}},
			TaskReader:        &mockTaskReader{},
			UserReader:        &mockUserReader{},
			PDFGenerator:      &mockPDFGenerator{},
			AIReportGenerator: &mockAIReportGenerator{text: "summary"},
		},
		nil,
	)

	_, err := useCase.Execute(applicationpkg.GenerateWeeklyReportsInput{WorkspaceID: 1, WeekID: 2})
	if !errors.Is(err, reportsDomain.ErrReportNoAssignmentsFound) {
		t.Fatalf("expected ErrReportNoAssignmentsFound, got %v", err)
	}
}

func TestGenerateWeeklyReportsRejectsWhenAssignmentHasNoTasksForWeek(t *testing.T) {
	useCase := applicationpkg.NewGenerateWeeklyReports(
		applicationpkg.GenerateWeeklyReportsDependencies{
			ReportRepo:      newMockReportRepository(),
			WorkspaceReader: &mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "WS"}},
			WeekReader: &mockWeekReader{week: &weeksDomain.Week{
				ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate,
			}},
			AssignmentReader:  &mockAssignmentReader{assignments: []assignmentsDomain.Assignment{{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 4}}},
			TaskReader:        &mockTaskReader{tasks: []tasksDomain.Task{}},
			UserReader:        &mockUserReader{user: &usersDomain.User{ID: 2, Name: "Ana"}},
			PDFGenerator:      &mockPDFGenerator{},
			AIReportGenerator: &mockAIReportGenerator{text: "summary"},
		},
		nil,
	)

	_, err := useCase.Execute(applicationpkg.GenerateWeeklyReportsInput{WorkspaceID: 1, WeekID: 2})
	if !errors.Is(err, reportsDomain.ErrReportNoTasksFoundForWeek) {
		t.Fatalf("expected ErrReportNoTasksFoundForWeek, got %v", err)
	}
}

func TestGenerateWeeklyReportsSuccess(t *testing.T) {
	repo := newMockReportRepository()
	fileStorage := &mockReportFileStorage{}
	useCase := applicationpkg.NewGenerateWeeklyReports(
		applicationpkg.GenerateWeeklyReportsDependencies{
			ReportRepo:      repo,
			WorkspaceReader: &mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "Algorithms"}},
			WeekReader: &mockWeekReader{week: &weeksDomain.Week{
				ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate,
			}},
			AssignmentReader: &mockAssignmentReader{assignments: []assignmentsDomain.Assignment{{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 4}}},
			TaskReader: &mockTaskReader{tasks: []tasksDomain.Task{
				{ID: 5, AssignmentID: 10, Title: "Prepare lab", Description: "Slides", Status: tasksDomain.TaskStatusFinalizado, SpentHours: 2, WeekStartDate: time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)},
				{ID: 6, AssignmentID: 10, Title: "Support class", Description: "Questions", Status: tasksDomain.TaskStatusAbierto, SpentHours: 3, WeekStartDate: time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)},
			}},
			UserReader:        &mockUserReader{user: &usersDomain.User{ID: 2, Name: "Ana Gomez"}},
			PDFGenerator:      &mockPDFGenerator{},
			AIReportGenerator: &mockAIReportGenerator{text: "Weekly summary"},
			ReportFileStorage: fileStorage,
		},
		&applicationpkg.GenerateWeeklyReportsOptions{ReportsStorageDir: t.TempDir(), ReportsGCSPrefix: "reports"},
	)

	output, err := useCase.Execute(applicationpkg.GenerateWeeklyReportsInput{WorkspaceID: 1, WeekID: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(output.Reports))
	}
	if output.Reports[0].FilePath != "reports/workspace_1_week_2_assignment_10.pdf" {
		t.Fatalf("expected normalized file path, got %q", output.Reports[0].FilePath)
	}
	if fileStorage.objectName != "reports/workspace_1_week_2_assignment_10.pdf" {
		t.Fatalf("expected uploaded object name, got %q", fileStorage.objectName)
	}
	if len(fileStorage.uploaded) == 0 {
		t.Fatal("expected uploaded file bytes")
	}
}

func TestGenerateWeeklyReportsRejectsAIFailure(t *testing.T) {
	useCase := applicationpkg.NewGenerateWeeklyReports(
		applicationpkg.GenerateWeeklyReportsDependencies{
			ReportRepo:      newMockReportRepository(),
			WorkspaceReader: &mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "Algorithms"}},
			WeekReader: &mockWeekReader{week: &weeksDomain.Week{
				ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate,
			}},
			AssignmentReader: &mockAssignmentReader{assignments: []assignmentsDomain.Assignment{{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 4}}},
			TaskReader: &mockTaskReader{tasks: []tasksDomain.Task{
				{ID: 5, AssignmentID: 10, Title: "Prepare lab", Description: "Slides", Status: tasksDomain.TaskStatusFinalizado, SpentHours: 2},
			}},
			UserReader:        &mockUserReader{user: &usersDomain.User{ID: 2, Name: "Ana Gomez"}},
			PDFGenerator:      &mockPDFGenerator{},
			AIReportGenerator: &mockAIReportGenerator{err: errors.New("ai failed")},
		},
		&applicationpkg.GenerateWeeklyReportsOptions{ReportsStorageDir: t.TempDir()},
	)

	_, err := useCase.Execute(applicationpkg.GenerateWeeklyReportsInput{WorkspaceID: 1, WeekID: 2})
	if !errors.Is(err, reportsDomain.ErrReportAIGenerationFailed) {
		t.Fatalf("expected ErrReportAIGenerationFailed, got %v", err)
	}
}
