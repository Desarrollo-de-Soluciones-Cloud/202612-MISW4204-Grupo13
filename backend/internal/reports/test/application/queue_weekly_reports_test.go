package application_test

import (
	applicationpkg "backend/internal/reports/application"
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
	"time"
)

func TestQueueWeeklyReportsRejectsMissingWorkspaceID(t *testing.T) {
	useCase := applicationpkg.NewQueueWeeklyReports(
		&mockWorkspaceReader{},
		&mockWeekReader{},
		&mockAssignmentReader{},
		&mockTaskReader{},
		&mockReportJobPublisher{},
	)

	_, err := useCase.Execute(applicationpkg.QueueWeeklyReportsInput{WeekID: 2})
	if !errors.Is(err, reportsDomain.ErrReportWorkspaceIDRequired) {
		t.Fatalf("expected ErrReportWorkspaceIDRequired, got %v", err)
	}
}

func TestQueueWeeklyReportsRejectsWhenNoAssignmentsAreReportable(t *testing.T) {
	useCase := applicationpkg.NewQueueWeeklyReports(
		&mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "Algorithms"}},
		&mockWeekReader{week: &weeksDomain.Week{ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate}},
		&mockAssignmentReader{assignments: []assignmentsDomain.Assignment{
			{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.AssignmentRole("professor"), WeeklyHours: 4},
		}},
		&mockTaskReader{},
		&mockReportJobPublisher{},
	)

	_, err := useCase.Execute(applicationpkg.QueueWeeklyReportsInput{WorkspaceID: 1, WeekID: 2})
	if !errors.Is(err, reportsDomain.ErrReportNoAssignmentsFound) {
		t.Fatalf("expected ErrReportNoAssignmentsFound, got %v", err)
	}
}

func TestQueueWeeklyReportsRejectsWhenNoTasksExistForReportableAssignments(t *testing.T) {
	useCase := applicationpkg.NewQueueWeeklyReports(
		&mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "Algorithms"}},
		&mockWeekReader{week: &weeksDomain.Week{ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate}},
		&mockAssignmentReader{assignments: []assignmentsDomain.Assignment{
			{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 4},
		}},
		&mockTaskReader{},
		&mockReportJobPublisher{},
	)

	_, err := useCase.Execute(applicationpkg.QueueWeeklyReportsInput{WorkspaceID: 1, WeekID: 2})
	if !errors.Is(err, reportsDomain.ErrReportNoTasksFoundForWeek) {
		t.Fatalf("expected ErrReportNoTasksFoundForWeek, got %v", err)
	}
}

func TestQueueWeeklyReportsPublishesOnlyAssignmentsWithTasks(t *testing.T) {
	publisher := &mockReportJobPublisher{}
	useCase := applicationpkg.NewQueueWeeklyReports(
		&mockWorkspaceReader{workspace: &workspacesDomain.Workspace{ID: 1, Name: "Algorithms"}},
		&mockWeekReader{week: &weeksDomain.Week{ID: 2, Number: 3, InitialDate: testWeekInitialDate, FinalDate: testWeekFinalDate}},
		&mockAssignmentReader{assignments: []assignmentsDomain.Assignment{
			{ID: 10, UserID: 2, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 4},
			{ID: 11, UserID: 3, WorkspaceID: 1, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 6},
		}},
		&mockTaskReader{tasks: []tasksDomain.Task{
			{ID: 5, AssignmentID: 11, Title: "Support class", Status: tasksDomain.TaskStatusFinalizado, SpentHours: 3, WeekStartDate: time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)},
		}},
		publisher,
	)

	output, err := useCase.Execute(applicationpkg.QueueWeeklyReportsInput{WorkspaceID: 1, WeekID: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.QueuedCount != 1 {
		t.Fatalf("expected queued count 1, got %d", output.QueuedCount)
	}
	if len(publisher.jobs) != 1 {
		t.Fatalf("expected 1 published job, got %d", len(publisher.jobs))
	}
	if publisher.jobs[0].AssignmentID != 11 || publisher.jobs[0].UserID != 3 {
		t.Fatalf("unexpected published job: %+v", publisher.jobs[0])
	}
}
