package application_test

import (
	usersDomain "backend/internal/users/domain"
	workspacesApplication "backend/internal/workspaces/application"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
)

func TestCloseWorkspaceSuccess(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{
		workspace: &workspacesDomain.Workspace{
			ID:           1,
			PeriodID:     2,
			UserID:       10,
			Name:         "Algorithms",
			Type:         workspacesDomain.CourseType,
			InitialDate:  "2026-06-01",
			FinalDate:    "2026-06-30",
			Observations: "obs",
			State:        workspacesDomain.ActiveState,
		},
	}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewCloseWorkspace(workspaceRepo, userRepo)
	output, err := uc.Execute(workspacesApplication.CloseWorkspaceInput{ID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.State != string(workspacesDomain.ClosedState) {
		t.Fatalf("expected closed state, got %q", output.State)
	}
}

func TestCloseWorkspaceRejectsNonProfessorOwner(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{
		workspace: &workspacesDomain.Workspace{ID: 1, UserID: 10, State: workspacesDomain.ActiveState},
	}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleMonitor}}

	uc := workspacesApplication.NewCloseWorkspace(workspaceRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CloseWorkspaceInput{ID: 1})
	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserNotProfessor) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceUserNotProfessor, err)
	}
}
