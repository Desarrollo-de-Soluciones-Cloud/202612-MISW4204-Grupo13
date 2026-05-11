package application_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	usersDomain "backend/internal/users/domain"
	workspacesApplication "backend/internal/workspaces/application"
	workspacesDomain "backend/internal/workspaces/domain"
	"testing"
)

func TestListWorkspaceMonitorsAndAssistantsEmpty(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	assignmentRepo := &assignmentRepoStub{}
	userRepo := &userRepoStub{}

	uc := workspacesApplication.NewListWorkspaceMonitorsAndAssistants(workspaceRepo, assignmentRepo, userRepo)
	output, err := uc.Execute(workspacesApplication.ListWorkspaceMonitorsAndAssistantsInput{ProfessorID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Monitors) != 0 || len(output.Assistants) != 0 {
		t.Fatal("expected empty monitor/assistant lists")
	}
}

func TestListWorkspaceMonitorsAndAssistantsSuccess(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{
		workspaces: []workspacesDomain.Workspace{
			{ID: 5, UserID: 10, Name: "Algorithms", Type: workspacesDomain.CourseType, State: workspacesDomain.ActiveState},
		},
	}
	assignmentRepo := &assignmentRepoStub{
		assignments: []assignmentsDomain.Assignment{
			{ID: 1, WorkspaceID: 5, UserID: 20, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 6},
			{ID: 2, WorkspaceID: 5, UserID: 30, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 8},
		},
	}
	userRepo := &userRepoStub{
		user: &usersDomain.User{ID: 30, Name: "Assistant", Email: "assistant@example.com", GlobalRole: usersDomain.RoleAssistant},
	}

	uc := workspacesApplication.NewListWorkspaceMonitorsAndAssistants(workspaceRepo, assignmentRepo, userRepo)
	output, err := uc.Execute(workspacesApplication.ListWorkspaceMonitorsAndAssistantsInput{ProfessorID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Assistants) != 1 {
		t.Fatalf("expected 1 assistant, got %d", len(output.Assistants))
	}
}
