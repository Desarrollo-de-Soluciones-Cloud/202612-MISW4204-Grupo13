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
	testDateLayout       = "2006-01-02"
	testWorkspaceStart   = "2026-06-01"
	testWorkspaceEnd     = "2026-06-30"
	testExpectedGotMsg   = "expected %v, got %v"
)

func TestCreateWorkspaceRejectsClosedPeriod(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ClosedPeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format(testDateLayout)}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: testWorkspaceStart, FinalDate: testWorkspaceEnd, Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspacePeriodClosed) {
		t.Fatalf(testExpectedGotMsg, workspacesDomain.ErrWorkspacePeriodClosed, err)
	}
}

func TestCreateWorkspaceRejectsAfterInscriptionDeadline(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, -1).Format(testDateLayout)}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: testWorkspaceStart, FinalDate: testWorkspaceEnd, Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceInscriptionClosed) {
		t.Fatalf(testExpectedGotMsg, workspacesDomain.ErrWorkspaceInscriptionClosed, err)
	}
}

func TestCreateWorkspaceRejectsMissingProfessor(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InitialDate: testWorkspaceStart, FinalDate: testWorkspaceEnd, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format(testDateLayout)}}
	userRepo := &userRepoStub{err: usersDomain.ErrUserNotFound}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: testWorkspaceStart, FinalDate: testWorkspaceEnd, Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserNotFound) {
		t.Fatalf(testExpectedGotMsg, workspacesDomain.ErrWorkspaceUserNotFound, err)
	}
}

func TestCreateWorkspaceRejectsNonProfessorOwner(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InitialDate: testWorkspaceStart, FinalDate: testWorkspaceEnd, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format(testDateLayout)}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleMonitor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: testWorkspaceStart, FinalDate: testWorkspaceEnd, Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserNotProfessor) {
		t.Fatalf(testExpectedGotMsg, workspacesDomain.ErrWorkspaceUserNotProfessor, err)
	}
}

func TestCreateWorkspaceSuccess(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{
		period: &periodsDomain.Period{
			ID:                   1,
			PeriodState:          periodsDomain.ActivePeriod,
			InitialDate:          testWorkspaceStart,
			FinalDate:            testWorkspaceEnd,
			InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format(testDateLayout),
		},
	}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	output, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{
		PeriodID:     1,
		UserID:       10,
		Name:         "Algorithms",
		Type:         "course",
		InitialDate:  "2026-06-02",
		FinalDate:    testWorkspaceEnd,
		Observations: "obs",
		State:        "active",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.ID == 0 {
		t.Fatal("expected persisted workspace id")
	}
}
