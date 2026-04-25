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

type workspaceRepoStub struct {
	workspace *workspacesDomain.Workspace
	createErr error
	updateErr error
}

func (r *workspaceRepoStub) Create(workspace *workspacesDomain.Workspace) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.workspace = workspace
	if r.workspace.ID == 0 {
		r.workspace.ID = 1
	}
	return nil
}

func (r *workspaceRepoStub) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	if r.workspace == nil || r.workspace.ID != id {
		return nil, workspacesDomain.ErrWorkspaceNotFound
	}
	return r.workspace, nil
}

func (r *workspaceRepoStub) FindAll() ([]workspacesDomain.Workspace, error) { return nil, nil }
func (r *workspaceRepoStub) FindByPeriodID(periodID uint) ([]workspacesDomain.Workspace, error) {
	return nil, nil
}
func (r *workspaceRepoStub) Update(workspace *workspacesDomain.Workspace) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	r.workspace = workspace
	return nil
}
func (r *workspaceRepoStub) Delete(id uint) error { return nil }

type periodRepoStub struct {
	period *periodsDomain.Period
	err    error
}

func (r *periodRepoStub) Create(period *periodsDomain.Period) error { return nil }
func (r *periodRepoStub) FindByID(id uint) (*periodsDomain.Period, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.period == nil || r.period.ID != id {
		return nil, periodsDomain.ErrPeriodNotFound
	}
	return r.period, nil
}
func (r *periodRepoStub) FindByName(name string) (*periodsDomain.Period, error) { return nil, periodsDomain.ErrPeriodNotFound }
func (r *periodRepoStub) FindAll() ([]periodsDomain.Period, error) { return nil, nil }
func (r *periodRepoStub) FindAllByState(state periodsDomain.PeriodState) ([]periodsDomain.Period, error) {
	return nil, nil
}
func (r *periodRepoStub) Update(period *periodsDomain.Period) error { return nil }
func (r *periodRepoStub) Delete(id uint) error { return nil }

type userRepoStub struct {
	user *usersDomain.User
	err  error
}

func (r *userRepoStub) Create(user *usersDomain.User) error { return nil }
func (r *userRepoStub) FindByID(id uint) (*usersDomain.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.user == nil || r.user.ID != id {
		return nil, usersDomain.ErrUserNotFound
	}
	return r.user, nil
}
func (r *userRepoStub) FindByEmail(email string) (*usersDomain.User, error) { return nil, usersDomain.ErrUserNotFound }
func (r *userRepoStub) FindAll() ([]usersDomain.User, error) { return nil, nil }
func (r *userRepoStub) FindAllByRole(role usersDomain.UserRole) ([]usersDomain.User, error) { return nil, nil }
func (r *userRepoStub) Update(user *usersDomain.User) error { return nil }

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
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}}
	userRepo := &userRepoStub{err: usersDomain.ErrUserNotFound}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserNotFound) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceUserNotFound, err)
	}
}

func TestCreateWorkspaceRejectsNonProfessorOwner(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleMonitor}}

	uc := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.CreateWorkspaceInput{PeriodID: 1, UserID: 10, Name: "WS", Type: "course", InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserNotProfessor) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceUserNotProfessor, err)
	}
}

func TestUpdateWorkspaceRejectsClosedWorkspace(t *testing.T) {
	workspaceRepo := &workspaceRepoStub{workspace: &workspacesDomain.Workspace{ID: 1, PeriodID: 1, UserID: 10, Name: "WS", Type: workspacesDomain.CourseType, InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: workspacesDomain.ClosedState}}
	periodRepo := &periodRepoStub{period: &periodsDomain.Period{ID: 1, PeriodState: periodsDomain.ActivePeriod, InscriptionFinalDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02")}}
	userRepo := &userRepoStub{user: &usersDomain.User{ID: 10, GlobalRole: usersDomain.RoleProfessor}}

	uc := workspacesApplication.NewUpdateWorkspace(workspaceRepo, periodRepo, userRepo)
	_, err := uc.Execute(workspacesApplication.UpdateWorkspaceInput{ID: 1, PeriodID: 1, Name: "Updated", Type: "course", InitialDate: "2026-06-01", FinalDate: "2026-06-30", Observations: "obs", State: "active"})

	if !errors.Is(err, workspacesDomain.ErrWorkspaceClosedUpdateForbidden) {
		t.Fatalf("expected %v, got %v", workspacesDomain.ErrWorkspaceClosedUpdateForbidden, err)
	}
}