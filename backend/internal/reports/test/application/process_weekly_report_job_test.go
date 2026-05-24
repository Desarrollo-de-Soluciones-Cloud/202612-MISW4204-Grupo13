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

func TestProcessWeeklyReportJobRejectsMissingAssignmentID(t *testing.T) {
	useCase := applicationpkg.NewProcessWeeklyReportJob(
		applicationpkg.ProcessWeeklyReportJobDependencies{
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

	_, err := useCase.Execute(applicationpkg.WeeklyReportJobMessage{WorkspaceID: 1, WeekID: 2, UserID: 4})
	if !errors.Is(err, reportsDomain.ErrReportAssignmentIDRequired) {
		t.Fatalf("expected ErrReportAssignmentIDRequired, got %v", err)
	}
}

func TestProcessWeeklyReportJobRejectsNonReportableAssignment(t *testing.T) {
	useCase := applicationpkg.NewProcessWeeklyReportJob(
		applicationpkg.ProcessWeeklyReportJobDependencies{
			ReportRepo:       newMockReportRepository(),
			WorkspaceReader:  &mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "Algorithms"}},
			WeekReader:       &mockWeekReader{week: &weeksDomain.Week{ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate}},
			AssignmentReader: &mockAssignmentReader{assignments: []assignmentsDomain.Assignment{{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.AssignmentRole("professor"), WeeklyHours: 4}}},
			TaskReader:       &mockTaskReader{},
			UserReader:       &mockUserReader{},
			PDFGenerator:     &mockPDFGenerator{},
			AIReportGenerator: &mockAIReportGenerator{},
		},
		nil,
	)

	_, err := useCase.Execute(applicationpkg.WeeklyReportJobMessage{WorkspaceID: 1, WeekID: 2, AssignmentID: 10, UserID: 2})
	if !errors.Is(err, reportsDomain.ErrReportNoAssignmentsFound) {
		t.Fatalf("expected ErrReportNoAssignmentsFound, got %v", err)
	}
}

func TestProcessWeeklyReportJobGeneratesReportForMatchingAssignment(t *testing.T) {
	repo := newMockReportRepository()
	fileStorage := &mockReportFileStorage{}
	useCase := applicationpkg.NewProcessWeeklyReportJob(
		applicationpkg.ProcessWeeklyReportJobDependencies{
			ReportRepo:      repo,
			WorkspaceReader: &mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "Algorithms"}},
			WeekReader:      &mockWeekReader{week: &weeksDomain.Week{ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate}},
			AssignmentReader: &mockAssignmentReader{assignments: []assignmentsDomain.Assignment{
				{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 4},
			}},
			TaskReader: &mockTaskReader{tasks: []tasksDomain.Task{
				{ID: 5, AssignmentID: 10, Title: "Prepare lab", Description: "Slides", Status: tasksDomain.TaskStatusFinalizado, SpentHours: 2, WeekStartDate: time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)},
			}},
			UserReader:        &mockUserReader{user: &usersDomain.User{ID: 2, Name: "Ana Gomez"}},
			PDFGenerator:      &mockPDFGenerator{},
			AIReportGenerator: &mockAIReportGenerator{text: "summary"},
			ReportFileStorage: fileStorage,
		},
		&applicationpkg.GenerateWeeklyReportsOptions{ReportsGCSPrefix: "reports"},
	)

	output, err := useCase.Execute(applicationpkg.WeeklyReportJobMessage{WorkspaceID: 1, WeekID: 2, AssignmentID: 10, UserID: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.AssignmentID != 10 || output.UserID != 2 {
		t.Fatalf("unexpected output: %+v", output)
	}
	if len(repo.reports) != 1 {
		t.Fatalf("expected one persisted report, got %d", len(repo.reports))
	}
	if fileStorage.objectName == "" {
		t.Fatalf("expected report to be uploaded to storage")
	}
}
