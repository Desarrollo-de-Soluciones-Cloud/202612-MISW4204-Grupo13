package application

import (
	applicationpkg "backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	"errors"
	"testing"
	"time"
)

func TestDeleteTaskRejectsInactiveWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	task := seedTask(t, taskRepo, 1, 10, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	deleteTask := applicationpkg.NewDeleteTask(taskRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})

	err := deleteTask.Execute(applicationpkg.DeleteTaskInput{ID: task.ID})
	if !errors.Is(err, domain.ErrTaskDeleteForbidden) {
		t.Fatalf("expected ErrTaskDeleteForbidden, got %v", err)
	}
}
