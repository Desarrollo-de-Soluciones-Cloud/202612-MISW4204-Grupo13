package infrastructure_test

import (
	tasksInfrastructure "backend/internal/tasks/infrastructure"
	"testing"
)

func TestNewTaskRepository(t *testing.T) {
	repo := tasksInfrastructure.NewTaskRepository()
	if repo == nil {
		t.Fatalf("expected task repository")
	}
}
