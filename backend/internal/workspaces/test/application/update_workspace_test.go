package application_test

import (
	periodsDomain "backend/internal/periods/domain"
	usersDomain "backend/internal/users/domain"
	workspacesApplication "backend/internal/workspaces/application"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
	"time"
)

const (
	testWorkspaceDateLayout = "2006-01-02"
	testWorkspaceStartDate  = "2026-06-01"
	testWorkspaceEndDate    = "2026-06-30"
)

func TestUpdateWorkspaceRejectsClosedWorkspace(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{workspace: &workspacesDomain.Workspace{ID: 1, PeriodID: 1, UserID: 10, Name: "WS", Type: workspacesDomain.CourseType, InitialDate: testWorkspaceStartDate, FinalDate: testWorkspaceEndDate, Observations: "obs", State: workspacesDomain.ClosedState}}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format(testWorkspaceDateLayout)}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewUpdateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.UpdateWorkspaceInput{ID: 1, PeriodID: 1, Name: "Updated", Type: "course", InitialDate: testWorkspaceStartDate, FinalDate: testWorkspaceEndDate, Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceClosedUpdateForbidden) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceClosedUpdateForbidden, err)
	}
}

func TestUpdateWorkspaceSuccess(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{
		workspace: &workspacesDomain.Workspace{
			ID:           1,
			PeriodID:     1,
			UserID:       10,
			Name:         "WS",
			Type:         workspacesDomain.CourseType,
			InitialDate:  testWorkspaceStartDate,
			FinalDate:    testWorkspaceEndDate,
			Observations: "obs",
			State:        workspacesDomain.ActiveState,
		},
	}
	periodRepo := &periodRepoStub{
		period: &periodsDomain.Period{
			ID:                   1,
			PeriodState:          periodsDomain.ActivePeriod,
			InitialDate:          testWorkspaceStartDate,
			FinalDate:            testWorkspaceEndDate,
			InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format(testWorkspaceDateLayout),
		},
	}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewUpdateWorkspace(workspaceRepo, periodRepo, userRepo)
	output, err := uc.Execute(workspacesApplication.UpdateWorkspaceInput{
		ID:           1,
		PeriodID:     1,
		UserID:       10,
		Name:         "Updated",
		Type:         "course",
		InitialDate:  "2026-06-02",
		FinalDate:    "2026-06-29",
		Observations: "updated",
		State:        "active",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "Updated" {
		t.Fatalf("expected updated name, got %q", output.Name)
	}
}
