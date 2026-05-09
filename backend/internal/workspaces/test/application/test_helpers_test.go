package application_test

import (
	periodsDomain "backend/internal/periods/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
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

func (r *workspaceRepoStub) FindByUserID(userID uint) ([]workspacesDomain.Workspace, error) {
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

func (r *periodRepoStub) FindByName(name string) (*periodsDomain.Period, error) {
	return nil, periodsDomain.ErrPeriodNotFound
}

func (r *periodRepoStub) FindAll() ([]periodsDomain.Period, error) { return nil, nil }

func (r *periodRepoStub) FindAllByState(state periodsDomain.PeriodState) ([]periodsDomain.Period, error) {
	return nil, nil
}

func (r *periodRepoStub) Update(period *periodsDomain.Period) error { return nil }
func (r *periodRepoStub) Delete(id uint) error                      { return nil }

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

func (r *userRepoStub) FindByEmail(email string) (*usersDomain.User, error) {
	return nil, usersDomain.ErrUserNotFound
}

func (r *userRepoStub) FindAll() ([]usersDomain.User, error) { return nil, nil }

func (r *userRepoStub) FindAllByRole(role usersDomain.UserRole) ([]usersDomain.User, error) {
	return nil, nil
}

func (r *userRepoStub) Update(user *usersDomain.User) error { return nil }
