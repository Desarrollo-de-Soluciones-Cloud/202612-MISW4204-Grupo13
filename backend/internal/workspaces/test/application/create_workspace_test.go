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

func TestCreateWorkspaceRejectsClosedPeriod(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ClosedPeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspacePeriodClosed) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspacePeriodClosed, err)
	}
}

func TestCreateWorkspaceRejectsAfterInscriptionDeadline(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, -1).Format("2006-01-02")}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceInscriptionClosed) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceInscriptionClosed, err)
	}
}

func TestCreateWorkspaceRejectsMissingProfessor(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InitialDate: "2026-06-01", FinalDate: "2026-06-30", InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}}
	userRepo := &userRepoStub{err: usersDomain.ErrUserNotFound}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserNotFound) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceUserNotFound, err)
	}
}

func TestCreateWorkspaceRejectsNonProfessorOwner(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InitialDate: "2026-06-01", FinalDate: "2026-06-30", InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleMonitor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserNotProfessor) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceUserNotProfessor, err)
	}
}
