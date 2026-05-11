package application

import (
	applicationpkg "backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	"errors"
	"testing"
	"time"
)

func TestUpdateTaskRejectsLateTask(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	task := seedTask(t, taskRepo, 1, 10, true, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})

	_, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:            task.ID,
		AssignmentID:  10,
		Title:         "Prepare class 2",
		Description:   "Review slides",
		Status:        domain.TaskStatusFinalizado,
		SpentHours:    3,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrTaskLateUpdateForbidden) {
		t.Fatalf("expected ErrTaskLateUpdateForbidden, got %v", err)
	}
}

func TestUpdateTaskAllowsChangingStatusDuringActiveWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	task := seedTask(t, taskRepo, 1, 10, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, func() time.Time {
		return time.Date(2026, 4, 9, 9, 0, 0, 0, time.UTC)
	})

	output, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:            task.ID,
		AssignmentID:  10,
		Title:         "Prepare class",
		Description:   "Review slides",
		Status:        domain.TaskStatusFinalizado,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Status != domain.TaskStatusFinalizado {
		t.Fatalf("expected updated status %q, got %q", domain.TaskStatusFinalizado, output.Status)
	}
}
