package application

import (
	applicationpkg "backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
	"time"
)

const (
	testCreateTaskTitle       = "Prepare class"
	testCreateTaskDescription = "Review slides"
)

func TestCreateTaskRejectsInactiveWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 3)
	seedWorkspace(workspaceRepo, 5, workspacesDomain.ActiveState)
	now := func() time.Time { return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC) }
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, now)

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID:  10,
		Title:         testCreateTaskTitle,
		Description:   testCreateTaskDescription,
		Status:        domain.TaskStatusAbierto,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, weeksDomain.ErrWeekNotFound) {
		t.Fatalf("expected ErrWeekNotFound, got %v", err)
	}
}

func TestCreateTaskRejectsMissingAssignment(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, time.Now)

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID:  10,
		Title:         testCreateTaskTitle,
		Description:   testCreateTaskDescription,
		Status:        domain.TaskStatusAbierto,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrTaskAssignmentNotFound) {
		t.Fatalf("expected ErrTaskAssignmentNotFound, got %v", err)
	}
}

func TestCreateTaskRejectsLegacyEnglishStatus(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 3)
	seedWorkspace(workspaceRepo, 5, workspacesDomain.ActiveState)
	seedWeek(weekRepo, 1, "2026-04-06")
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 9, 9, 0, 0, 0, time.UTC)
	})

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID:  10,
		Title:         testCreateTaskTitle,
		Description:   testCreateTaskDescription,
		Status:        domain.TaskStatus("open"),
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrTaskStatusInvalid) {
		t.Fatalf("expected ErrTaskStatusInvalid, got %v", err)
	}
}

func TestCreateTaskRejectsClosedWorkspace(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 3)
	seedWorkspace(workspaceRepo, 5, workspacesDomain.ClosedState)
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 9, 9, 0, 0, 0, time.UTC)
	})

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID:  10,
		Title:         testCreateTaskTitle,
		Description:   testCreateTaskDescription,
		Status:        domain.TaskStatusAbierto,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrTaskWorkspaceClosed) {
		t.Fatalf("expected ErrTaskWorkspaceClosed, got %v", err)
	}
}
