package delivery_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	authDomain "backend/internal/auth/domain"
	periodsDomain "backend/internal/periods/domain"
	applicationpkg "backend/internal/workspaces/application"
	deliverypkg "backend/internal/workspaces/delivery"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type workspaceRepoStub struct {
	workspace  *workspacesDomain.Workspace
	workspaces []workspacesDomain.Workspace
}

type periodRepoStub struct {
	period *periodsDomain.Period
	err    error
}

type userRepoStub struct {
	user  *usersDomain.User
	users map[uint]*usersDomain.User
	err   error
}

type assignmentRepoStub struct {
	assignments []assignmentsDomain.Assignment
	err         error
}

func (r *workspaceRepoStub) Create(workspace *workspacesDomain.Workspace) error {
	r.workspace = workspace
	if r.workspace.ID == 0 {
		r.workspace.ID = 1
	}
	return nil
}

func (r *workspaceRepoStub) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	if r.workspace != nil && r.workspace.ID == id {
		return r.workspace, nil
	}
	for i := range r.workspaces {
		if r.workspaces[i].ID == id {
			return &r.workspaces[i], nil
		}
	}
	return nil, workspacesDomain.ErrWorkspaceNotFound
}

func (r *workspaceRepoStub) FindAll() ([]workspacesDomain.Workspace, error) {
	return r.workspaces, nil
}

func (r *workspaceRepoStub) FindByPeriodID(periodID uint) ([]workspacesDomain.Workspace, error) {
	return r.workspaces, nil
}

func (r *workspaceRepoStub) FindByUserID(userID uint) ([]workspacesDomain.Workspace, error) {
	return r.workspaces, nil
}

func (r *workspaceRepoStub) Update(workspace *workspacesDomain.Workspace) error {
	r.workspace = workspace
	return nil
}

func (r *workspaceRepoStub) Delete(id uint) error {
	return nil
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

func (r *userRepoStub) Create(user *usersDomain.User) error { return nil }
func (r *userRepoStub) FindByID(id uint) (*usersDomain.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.users != nil {
		if user, ok := r.users[id]; ok {
			return user, nil
		}
	}
	if r.user != nil && r.user.ID == id {
		return r.user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}
func (r *userRepoStub) FindByEmail(email string) (*usersDomain.User, error) {
	return nil, usersDomain.ErrUserNotFound
}
func (r *userRepoStub) FindAll() ([]usersDomain.User, error) { return nil, nil }
func (r *userRepoStub) FindAllByRole(role usersDomain.UserRole) ([]usersDomain.User, error) {
	return nil, nil
}
func (r *userRepoStub) Update(user *usersDomain.User) error { return nil }

func (r *assignmentRepoStub) Create(assignment *assignmentsDomain.Assignment) error { return nil }
func (r *assignmentRepoStub) FindByID(id uint) (*assignmentsDomain.Assignment, error) {
	return nil, assignmentsDomain.ErrAssignmentNotFound
}
func (r *assignmentRepoStub) FindAll() ([]assignmentsDomain.Assignment, error) { return nil, nil }
func (r *assignmentRepoStub) FindAllByUserID(userID uint) ([]assignmentsDomain.Assignment, error) {
	return nil, nil
}
func (r *assignmentRepoStub) FindByWorkspaceUserID(workspaceUserID uint) ([]assignmentsDomain.Assignment, error) {
	return nil, nil
}
func (r *assignmentRepoStub) FindByWorkspaceIDsAndRoles(workspaceIDs []uint, roles []assignmentsDomain.AssignmentRole) ([]assignmentsDomain.Assignment, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.assignments, nil
}
func (r *assignmentRepoStub) SumWeeklyHoursByUserAndRole(userID uint, role assignmentsDomain.AssignmentRole) (int, error) {
	return 0, nil
}
func (r *assignmentRepoStub) CountAssignmentsByUserAndRole(userID uint, role assignmentsDomain.AssignmentRole) (int, error) {
	return 0, nil
}
func (r *assignmentRepoStub) Update(assignment *assignmentsDomain.Assignment) error { return nil }

func newWorkspaceHandlerForTest() *deliverypkg.WorkspaceHandler {
	workspaceRepo := &workspaceRepoStub{
		workspaces: []workspacesDomain.Workspace{
			{ID: 1, PeriodID: 1, UserID: 10, Name: "Algorithms", Type: workspacesDomain.CourseType, InitialDate: "2026-06-02", FinalDate: "2026-06-30", Observations: "obs", State: workspacesDomain.ActiveState},
			{ID: 2, PeriodID: 1, UserID: 99, Name: "AI Lab", Type: workspacesDomain.ProjectType, InitialDate: "2026-06-02", FinalDate: "2026-06-30", Observations: "obs", State: workspacesDomain.ActiveState},
		},
		workspace: &workspacesDomain.Workspace{ID: 1, PeriodID: 1, UserID: 10, Name: "Algorithms", Type: workspacesDomain.CourseType, InitialDate: "2026-06-02", FinalDate: "2026-06-30", Observations: "obs", State: workspacesDomain.ActiveState},
	}
	periodRepo := &periodRepoStub{
		period: &periodsDomain.Period{
			ID:                   1,
			PeriodState:          periodsDomain.ActivePeriod,
			InitialDate:          "2026-06-01",
			FinalDate:            "2026-06-30",
			InscriptionFinalDate: "2099-06-01",
		},
	}
	userRepo := &userRepoStub{
		user: &usersDomain.User{ID: 10, Name: "Prof", Email: "prof@example.com", GlobalRole: usersDomain.RoleProfessor},
		users: map[uint]*usersDomain.User{
			10: {ID: 10, Name: "Prof", Email: "prof@example.com", GlobalRole: usersDomain.RoleProfessor},
			20: {ID: 20, Name: "Monitor", Email: "monitor@example.com", GlobalRole: usersDomain.RoleMonitor},
			30: {ID: 30, Name: "Assistant", Email: "assistant@example.com", GlobalRole: usersDomain.RoleAssistant},
		},
	}
	assignmentRepo := &assignmentRepoStub{
		assignments: []assignmentsDomain.Assignment{
			{ID: 1, WorkspaceID: 1, UserID: 20, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 6},
			{ID: 2, WorkspaceID: 1, UserID: 30, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 8},
		},
	}

	return deliverypkg.NewWorkspaceHandler(deliverypkg.WorkspaceHandlerUseCases{
		CreateWorkspace:                    applicationpkg.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo),
		ListWorkspaces:                     applicationpkg.NewListWorkspaces(workspaceRepo),
		ListWorkspacesByPeriod:             applicationpkg.NewListWorkspacesByPeriod(workspaceRepo),
		GetWorkspaceByID:                   applicationpkg.NewGetWorkspaceByID(workspaceRepo),
		UpdateWorkspace:                    applicationpkg.NewUpdateWorkspace(workspaceRepo, periodRepo, userRepo),
		DeleteWorkspace:                    applicationpkg.NewDeleteWorkspace(workspaceRepo),
		CloseWorkspace:                     applicationpkg.NewCloseWorkspace(workspaceRepo, userRepo),
		ListWorkspaceMonitorsAndAssistants: applicationpkg.NewListWorkspaceMonitorsAndAssistants(workspaceRepo, assignmentRepo, userRepo),
	})
}

func authenticatedUser(id uint, role usersDomain.UserRole) authDomain.AuthenticatedUser {
	return authDomain.AuthenticatedUser{ID: id, GlobalRole: role}
}
