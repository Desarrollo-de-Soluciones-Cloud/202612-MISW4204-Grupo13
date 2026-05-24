package infrastructure

import "testing"

func TestNewWorkspaceRepository(t *testing.T) {
	repo := NewWorkspaceRepository()
	if repo == nil {
		t.Fatalf("expected workspace repository")
	}
}
