package infrastructure_test

import (
	sharedDB "backend/internal/shared/database/testsupport"
	tasksDomain "backend/internal/tasks/domain"
	tasksinfra "backend/internal/tasks/infrastructure"
	"errors"
	"testing"
	"time"
)

func TestTaskRepositorySQLiteCRUD(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &tasksDomain.Task{})
	repo := tasksinfra.NewTaskRepository()

	task, err := tasksDomain.NewTask(tasksDomain.TaskInput{
		UserID: 2, AssignmentID: 10, Title: "Task", Description: "Desc",
		Status: tasksDomain.TaskStatusAbierto, SpentHours: 2, Observations: "obs",
		WeekStartDate: time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected task, got %v", err)
	}
	if err := repo.Create(task); err != nil {
		t.Fatalf("expected create, got %v", err)
	}

	found, err := repo.FindByID(task.ID)
	if err != nil || found.Title != "Task" {
		t.Fatalf("expected find by id, got %v %#v", err, found)
	}

	all, err := repo.FindAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("expected 1 task, got %v %d", err, len(all))
	}
}

func TestTaskRepositorySQLiteDeleteNotFound(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &tasksDomain.Task{})
	repo := tasksinfra.NewTaskRepository()

	if err := repo.Delete(999); !errors.Is(err, tasksDomain.ErrTaskNotFound) {
		t.Fatalf("expected task not found, got %v", err)
	}
}
