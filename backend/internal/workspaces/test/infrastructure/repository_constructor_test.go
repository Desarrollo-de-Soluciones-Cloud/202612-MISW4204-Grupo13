package infrastructure_test

import (
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"
	"testing"
)

func TestNewWorkspaceRepository(t *testing.T) {
	repo := workspacesInfrastructure.NewWorkspaceRepository()
	if repo == nil {
		t.Fatalf("expected workspace repository")
	}
}
