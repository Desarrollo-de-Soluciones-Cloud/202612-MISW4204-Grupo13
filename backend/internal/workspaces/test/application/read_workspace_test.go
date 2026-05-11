package application_test

import (
	workspacesApplication "backend/internal/workspaces/application"
	workspacesDomain "backend/internal/workspaces/domain"
	"testing"
)

func TestGetWorkspaceByIDSuccess(t *testing.T) {
	repo := &workspaceRepoStub{
		workspace: &workspacesDomain.Workspace{
			ID:           1,
			PeriodID:     2,
			UserID:       3,
			Name:         "Algorithms",
			Type:         workspacesDomain.CourseType,
			InitialDate:  "2026-06-01",
			FinalDate:    "2026-06-30",
			Observations: "obs",
			State:        workspacesDomain.ActiveState,
		},
	}

	uc := workspacesApplication.NewGetWorkspaceByID(repo)
	output, err := uc.Execute(workspacesApplication.GetWorkspaceByIDInput{ID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "Algorithms" {
		t.Fatalf("expected workspace name, got %q", output.Name)
	}
}

func TestListWorkspacesSuccess(t *testing.T) {
	repo := &workspaceRepoStub{
		workspaces: []workspacesDomain.Workspace{
			{ID: 1, Name: "Algorithms", Type: workspacesDomain.CourseType, State: workspacesDomain.ActiveState},
			{ID: 2, Name: "AI Lab", Type: workspacesDomain.ProjectType, State: workspacesDomain.ActiveState},
		},
	}

	uc := workspacesApplication.NewListWorkspaces(repo)
	output, err := uc.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Workspaces) != 2 {
		t.Fatalf("expected 2 workspaces, got %d", len(output.Workspaces))
	}
}

func TestListWorkspacesByPeriodSuccess(t *testing.T) {
	repo := &workspaceRepoStub{
		workspaces: []workspacesDomain.Workspace{
			{ID: 1, PeriodID: 2, Name: "Algorithms", Type: workspacesDomain.CourseType, State: workspacesDomain.ActiveState},
		},
	}

	uc := workspacesApplication.NewListWorkspacesByPeriod(repo)
	output, err := uc.Execute(workspacesApplication.ListWorkspacesByPeriodInput{PeriodID: 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Workspaces) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(output.Workspaces))
	}
}
