package infrastructure

import "testing"

func TestNewAssignmentRepository(t *testing.T) {
	repo := NewAssignmentRepository()
	if repo == nil {
		t.Fatalf("expected assignment repository")
	}
}
