package application

import (
	applicationpkg "backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type mockAssignmentRepository struct {
	byID    map[uint]*domain.Assignment
	nextID  uint
	failErr error
}

type mockUserRepository struct {
	users map[uint]*usersDomain.User
}

type mockWorkspaceRepository struct {
	workspaces map[uint]*workspacesDomain.Workspace
}

const (
	errCreateAssignmentMsg          = "expected no error creating assignment: %v"
	errCreateAssistantAssignmentMsg = "expected no error creating assistant assignment, got %v"
	errCreateMonitorToUpdateMsg     = "expected no error creating monitor assignment to update, got %v"
)

func newMockAssignmentRepository() *mockAssignmentRepository {
	return &mockAssignmentRepository{
		byID:   make(map[uint]*domain.Assignment),
		nextID: 1,
	}
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: make(map[uint]*usersDomain.User)}
}

func newMockWorkspaceRepository() *mockWorkspaceRepository {
	return &mockWorkspaceRepository{workspaces: make(map[uint]*workspacesDomain.Workspace)}
}

func newCreateAssignmentWithDependencies(
	assignmentRepo *mockAssignmentRepository,
	userRepo *mockUserRepository,
	workspaceRepo *mockWorkspaceRepository,
) *applicationpkg.CreateAssignment {
	return applicationpkg.NewCreateAssignment(assignmentRepo).WithRepositories(userRepo, workspaceRepo)
}

func (m *mockAssignmentRepository) Create(assignment *domain.Assignment) error {
	if m.failErr != nil {
		return m.failErr
	}

	assignment.ID = m.nextID
	m.nextID++
	m.byID[assignment.ID] = assignment
	return nil
}

func (m *mockAssignmentRepository) FindByID(id uint) (*domain.Assignment, error) {
	if assignment, ok := m.byID[id]; ok {
		return assignment, nil
	}
	return nil, domain.ErrAssignmentNotFound
}

func (m *mockAssignmentRepository) FindAllByUserID(userID uint) ([]domain.Assignment, error) {
	result := make([]domain.Assignment, 0)
	for _, assignment := range m.byID {
		if assignment.UserID == userID {
			result = append(result, *assignment)
		}
	}
	return result, nil
}

func (m *mockAssignmentRepository) SumWeeklyHoursByUserAndRole(userID uint, role domain.AssignmentRole) (int, error) {
	total := 0
	for _, assignment := range m.byID {
		if assignment.UserID == userID && assignment.Role == role {
			total += assignment.WeeklyHours
		}
	}

	return total, nil
}

func (m *mockAssignmentRepository) CountAssignmentsByUserAndRole(userID uint, role domain.AssignmentRole) (int, error) {
	total := 0
	for _, assignment := range m.byID {
		if assignment.UserID == userID && assignment.Role == role {
			total++
		}
	}

	return total, nil
}

func (m *mockAssignmentRepository) Update(assignment *domain.Assignment) error {
	if _, ok := m.byID[assignment.ID]; !ok {
		return domain.ErrAssignmentNotFound
	}
	m.byID[assignment.ID] = assignment
	return nil
}

func (m *mockAssignmentRepository) FindAll() ([]domain.Assignment, error) {
	result := make([]domain.Assignment, 0, len(m.byID))
	for _, assignment := range m.byID {
		result = append(result, *assignment)
	}
	return result, nil
}

func (m *mockAssignmentRepository) FindByWorkspaceUserID(workspaceUserID uint) ([]domain.Assignment, error) {
	return nil, nil
}

func (m *mockAssignmentRepository) FindByWorkspaceIDsAndRoles(workspaceIDs []uint, roles []domain.AssignmentRole) ([]domain.Assignment, error) {
	result := make([]domain.Assignment, 0)
	roleMap := make(map[domain.AssignmentRole]bool)
	for _, role := range roles {
		roleMap[role] = true
	}

	workspaceMap := make(map[uint]bool)
	for _, wsID := range workspaceIDs {
		workspaceMap[wsID] = true
	}

	for _, assignment := range m.byID {
		if workspaceMap[assignment.WorkspaceID] && roleMap[assignment.Role] {
			result = append(result, *assignment)
		}
	}
	return result, nil
}

func (m *mockUserRepository) FindByID(id uint) (*usersDomain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}

func (m *mockWorkspaceRepository) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	if workspace, ok := m.workspaces[id]; ok {
		return workspace, nil
	}
	return nil, workspacesDomain.ErrWorkspaceNotFound
}

func (m *mockWorkspaceRepository) FindByUserID(userID uint) ([]workspacesDomain.Workspace, error) {
	result := make([]workspacesDomain.Workspace, 0)
	for _, workspace := range m.workspaces {
		if workspace.UserID == userID {
			result = append(result, *workspace)
		}
	}
	return result, nil
}
