package application_test

import (
	workspacesApplication "backend/internal/workspaces/application"
	"testing"
)

func TestDeleteWorkspaceSuccess(t *testing.T) {
	repo := &workspaceRepoStub{}
	uc := workspacesApplication.NewDeleteWorkspace(repo)

	if err := uc.Execute(workspacesApplication.DeleteWorkspaceInput{ID: 1}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
