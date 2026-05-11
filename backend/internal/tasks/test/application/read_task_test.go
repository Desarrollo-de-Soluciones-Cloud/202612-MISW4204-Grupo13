package application

import (
	applicationpkg "backend/internal/tasks/application"
	"testing"
	"time"
)

func TestListTasksReturnsAllTasks(t *testing.T) {
	taskRepo := newMockTaskRepository()
	seedTask(t, taskRepo, 1, 10, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	seedTask(t, taskRepo, 2, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))

	listTasks := applicationpkg.NewListTasks(taskRepo)
	output, err := listTasks.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(output.Tasks))
	}
}

func TestGetTaskByIDReturnsTask(t *testing.T) {
	taskRepo := newMockTaskRepository()
	task := seedTask(t, taskRepo, 2, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	getTask := applicationpkg.NewGetTaskByID(taskRepo)

	output, err := getTask.Execute(applicationpkg.GetTaskByIDInput{ID: task.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.ID != task.ID {
		t.Fatalf("expected task id %d, got %d", task.ID, output.ID)
	}
}
