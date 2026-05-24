package infrastructure

import "testing"

func TestNewTaskRepository(t *testing.T) {
	repo := NewTaskRepository()
	if repo == nil {
		t.Fatalf("expected task repository")
	}
}
