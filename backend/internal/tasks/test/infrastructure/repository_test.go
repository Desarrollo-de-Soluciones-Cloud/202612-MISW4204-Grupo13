package infrastructure

import (
	"backend/internal/shared/database"
	tasksDomain "backend/internal/tasks/domain"
	tasksInfrastructure "backend/internal/tasks/infrastructure"
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db
}

func TestTaskRepositoryCRUD(t *testing.T) {
	setupTestDB(t)

	repo := tasksInfrastructure.NewTaskRepository()
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	task, err := tasksDomain.NewTask(
		1,
		10,
		20,
		"Prepare class",
		"Review slides",
		tasksDomain.TaskStatusAbierto,
		2,
		"",
		[]tasksDomain.TaskAttachment{{Path: "docs/guide.pdf"}},
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected new task, got %v", err)
	}

	if err := repo.Create(task); err != nil {
		t.Fatalf("expected create task, got %v", err)
	}

	foundByID, err := repo.FindByID(task.ID)
	if err != nil {
		t.Fatalf("expected find by id, got %v", err)
	}
	if foundByID.AssignmentID != 10 || foundByID.WeekID != 20 {
		t.Fatalf("expected assignment id 10 and week id 20, got %d and %d", foundByID.AssignmentID, foundByID.WeekID)
	}
	if len(foundByID.Attachments) != 1 || foundByID.Attachments[0].Path != "docs/guide.pdf" {
		t.Fatalf("expected attachment to be preloaded, got %+v", foundByID.Attachments)
	}

	foundByID.Title = "Prepare class 2"
	foundByID.Attachments = []tasksDomain.TaskAttachment{{Path: "docs/updated.pdf"}}
	if err := repo.Update(foundByID); err != nil {
		t.Fatalf("expected update task, got %v", err)
	}

	updatedTask, err := repo.FindByID(task.ID)
	if err != nil {
		t.Fatalf("expected find updated task, got %v", err)
	}
	if updatedTask.Title != "Prepare class 2" {
		t.Fatalf("expected updated title, got %q", updatedTask.Title)
	}
	if len(updatedTask.Attachments) != 1 || updatedTask.Attachments[0].Path != "docs/updated.pdf" {
		t.Fatalf("expected replaced attachments, got %+v", updatedTask.Attachments)
	}

	allTasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("expected find all, got %v", err)
	}
	if len(allTasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(allTasks))
	}

	userTasks, err := repo.FindAllByUserID(1)
	if err != nil {
		t.Fatalf("expected find by user id, got %v", err)
	}
	if len(userTasks) != 1 {
		t.Fatalf("expected 1 user task, got %d", len(userTasks))
	}
}

func TestTaskRepositoryReturnsNotFound(t *testing.T) {
	setupTestDB(t)

	repo := tasksInfrastructure.NewTaskRepository()
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	_, err := repo.FindByID(999)
	if !errors.Is(err, tasksDomain.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskRepositoryDelete(t *testing.T) {
	setupTestDB(t)

	repo := tasksInfrastructure.NewTaskRepository()
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	task, err := tasksDomain.NewTask(
		1,
		10,
		20,
		"Prepare class",
		"Review slides",
		tasksDomain.TaskStatusAbierto,
		2,
		"",
		nil,
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected new task, got %v", err)
	}
	if err := repo.Create(task); err != nil {
		t.Fatalf("expected create task, got %v", err)
	}

	if err := repo.Delete(task.ID); err != nil {
		t.Fatalf("expected delete task, got %v", err)
	}
	if _, err := repo.FindByID(task.ID); !errors.Is(err, tasksDomain.ErrTaskNotFound) {
		t.Fatalf("expected task to be deleted, got %v", err)
	}
}

func TestTaskRepositoryNormalizesLegacyStatuses(t *testing.T) {
	setupTestDB(t)

	repo := tasksInfrastructure.NewTaskRepository()
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	legacyTask := &tasksDomain.Task{
		UserID:        1,
		AssignmentID:  10,
		WeekID:        20,
		Title:         "Prepare class",
		Description:   "Review slides",
		Status:        tasksDomain.TaskStatus("open"),
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		Late:          false,
	}
	if err := database.DB.Create(legacyTask).Error; err != nil {
		t.Fatalf("expected create legacy task, got %v", err)
	}

	if err := repo.NormalizeLegacyStatuses(); err != nil {
		t.Fatalf("expected legacy normalization, got %v", err)
	}

	foundByID, err := repo.FindByID(legacyTask.ID)
	if err != nil {
		t.Fatalf("expected find by id, got %v", err)
	}
	if foundByID.Status != tasksDomain.TaskStatusAbierto {
		t.Fatalf("expected normalized status %q, got %q", tasksDomain.TaskStatusAbierto, foundByID.Status)
	}
}
