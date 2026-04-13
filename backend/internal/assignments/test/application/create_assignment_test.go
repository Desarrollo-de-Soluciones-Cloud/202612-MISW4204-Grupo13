package application

import (
	applicationpkg "backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
)

//nolint:godox // TODO RF04: Agregar tests de delivery e infrastructure para fortalecer la cobertura del modulo.
// TODO RF05: Reforzar RF05 con pruebas de delivery e integracion/E2E para respuestas HTTP bloqueantes en create.

type MockAssignmentRepository struct {
	byID    map[uint]*domain.Assignment
	nextID  uint
	failErr error
}

type MockUserRepository struct {
	users map[uint]*usersDomain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[uint]*usersDomain.User)}
}

func (m *MockUserRepository) FindByID(id uint) (*usersDomain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}

type MockWorkspaceRepository struct {
	workspaces map[uint]*workspacesDomain.Workspace
}

func NewMockWorkspaceRepository() *MockWorkspaceRepository {
	return &MockWorkspaceRepository{workspaces: make(map[uint]*workspacesDomain.Workspace)}
}

func (m *MockWorkspaceRepository) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	if workspace, ok := m.workspaces[id]; ok {
		return workspace, nil
	}
	return nil, workspacesDomain.ErrWorkspaceNotFound
}

func newCreateAssignmentWithDependencies(
	assignmentRepo *MockAssignmentRepository,
	userRepo *MockUserRepository,
	workspaceRepo *MockWorkspaceRepository,
) *applicationpkg.CreateAssignment {
	return applicationpkg.NewCreateAssignment(assignmentRepo).WithRepositories(userRepo, workspaceRepo)
}

func NewMockAssignmentRepository() *MockAssignmentRepository {
	return &MockAssignmentRepository{
		byID:   make(map[uint]*domain.Assignment),
		nextID: 1,
	}
}

func (m *MockAssignmentRepository) Create(assignment *domain.Assignment) error {
	if m.failErr != nil {
		return m.failErr
	}

	assignment.ID = m.nextID
	m.nextID++
	m.byID[assignment.ID] = assignment
	return nil
}

func (m *MockAssignmentRepository) FindByID(id uint) (*domain.Assignment, error) {
	if assignment, ok := m.byID[id]; ok {
		return assignment, nil
	}
	return nil, domain.ErrAssignmentNotFound
}

func (m *MockAssignmentRepository) FindAllByUserID(userID uint) ([]domain.Assignment, error) {
	result := make([]domain.Assignment, 0)
	for _, assignment := range m.byID {
		if assignment.UserID == userID {
			result = append(result, *assignment)
		}
	}
	return result, nil
}

func (m *MockAssignmentRepository) SumWeeklyHoursByUserAndRole(userID uint, role domain.AssignmentRole) (int, error) {
	total := 0
	for _, assignment := range m.byID {
		if assignment.UserID == userID && assignment.Role == role {
			total += assignment.WeeklyHours
		}
	}

	return total, nil
}

func (m *MockAssignmentRepository) CountAssignmentsByUserAndRole(userID uint, role domain.AssignmentRole) (int, error) {
	total := 0
	for _, assignment := range m.byID {
		if assignment.UserID == userID && assignment.Role == role {
			total++
		}
	}

	return total, nil
}

func (m *MockAssignmentRepository) Update(assignment *domain.Assignment) error {
	if _, ok := m.byID[assignment.ID]; !ok {
		return domain.ErrAssignmentNotFound
	}
	m.byID[assignment.ID] = assignment
	return nil
}

func TestCreateAssignmentSuccess(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	userRepo := NewMockUserRepository()
	workspaceRepo := NewMockWorkspaceRepository()
	userRepo.users[1] = &usersDomain.User{ID: 1}
	workspaceRepo.workspaces[2] = &workspacesDomain.Workspace{ID: 2, State: workspacesDomain.ActiveState}
	createAssignment := newCreateAssignmentWithDependencies(mockRepo, userRepo, workspaceRepo)

	output, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.ID == 0 {
		t.Fatalf("expected generated id, got %d", output.ID)
	}
	if output.UserID != 1 {
		t.Fatalf("expected user id 1, got %d", output.UserID)
	}
	if output.WorkspaceID != 2 {
		t.Fatalf("expected workspace id 2, got %d", output.WorkspaceID)
	}
	if output.Role != domain.RoleAssistant {
		t.Fatalf("expected role %q, got %q", domain.RoleAssistant, output.Role)
	}
	if output.WeeklyHours != 10 {
		t.Fatalf("expected weekly hours 10, got %d", output.WeeklyHours)
	}
}

func TestCreateAssignmentReturnsErrorWhenUserDoesNotExist(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	userRepo := NewMockUserRepository()
	workspaceRepo := NewMockWorkspaceRepository()
	workspaceRepo.workspaces[2] = &workspacesDomain.Workspace{ID: 2, State: workspacesDomain.ActiveState}
	createAssignment := newCreateAssignmentWithDependencies(mockRepo, userRepo, workspaceRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      99,
		WorkspaceID: 2,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if !errors.Is(err, domain.ErrAssignmentUserNotFound) {
		t.Fatalf("expected ErrAssignmentUserNotFound, got %v", err)
	}
}

func TestCreateAssignmentReturnsErrorWhenWorkspaceDoesNotExist(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	userRepo := NewMockUserRepository()
	workspaceRepo := NewMockWorkspaceRepository()
	userRepo.users[1] = &usersDomain.User{ID: 1}
	createAssignment := newCreateAssignmentWithDependencies(mockRepo, userRepo, workspaceRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 99,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if !errors.Is(err, domain.ErrAssignmentWorkspaceNotFound) {
		t.Fatalf("expected ErrAssignmentWorkspaceNotFound, got %v", err)
	}
}

func TestCreateAssignmentReturnsErrorWhenWorkspaceIsClosed(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	userRepo := NewMockUserRepository()
	workspaceRepo := NewMockWorkspaceRepository()
	userRepo.users[1] = &usersDomain.User{ID: 1}
	workspaceRepo.workspaces[2] = &workspacesDomain.Workspace{ID: 2, State: workspacesDomain.ClosedState}
	createAssignment := newCreateAssignmentWithDependencies(mockRepo, userRepo, workspaceRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if !errors.Is(err, domain.ErrAssignmentWorkspaceClosed) {
		t.Fatalf("expected ErrAssignmentWorkspaceClosed, got %v", err)
	}
}

func TestCreateAssignmentAllowsIndependentAssignments(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	userRepo := NewMockUserRepository()
	workspaceRepo := NewMockWorkspaceRepository()
	userRepo.users[1] = &usersDomain.User{ID: 1}
	workspaceRepo.workspaces[10] = &workspacesDomain.Workspace{ID: 10, State: workspacesDomain.ActiveState}
	workspaceRepo.workspaces[11] = &workspacesDomain.Workspace{ID: 11, State: workspacesDomain.ActiveState}
	createAssignment := newCreateAssignmentWithDependencies(mockRepo, userRepo, workspaceRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 10,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf("expected no error creating first independent assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 11,
		Role:        domain.RoleAssistant,
		WeeklyHours: 8,
	})
	if err != nil {
		t.Fatalf("expected no error creating assignment in different workspace, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 10,
		Role:        domain.RoleMonitor,
		WeeklyHours: 4,
	})
	if err != nil {
		t.Fatalf("expected no error creating assignment with different role in same workspace, got %v", err)
	}
}

func TestCreateAssignmentBlocksExactDuplicate(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	userRepo := NewMockUserRepository()
	workspaceRepo := NewMockWorkspaceRepository()
	userRepo.users[1] = &usersDomain.User{ID: 1}
	workspaceRepo.workspaces[10] = &workspacesDomain.Workspace{ID: 10, State: workspacesDomain.ActiveState}
	createAssignment := newCreateAssignmentWithDependencies(mockRepo, userRepo, workspaceRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 10,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf("expected no error creating baseline assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 10,
		Role:        domain.RoleAssistant,
		WeeklyHours: 6,
	})
	if !errors.Is(err, domain.ErrAssignmentAlreadyExists) {
		t.Fatalf("expected ErrAssignmentAlreadyExists, got %v", err)
	}
}

func TestCreateAssignmentInvalidRole(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.AssignmentRole("invalid"),
		WeeklyHours: 10,
	})
	if !errors.Is(err, domain.ErrAssignmentRoleInvalid) {
		t.Fatalf("expected ErrAssignmentRoleInvalid, got %v", err)
	}
}

func TestCreateAssignmentInvalidWeeklyHours(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 0,
	})
	if !errors.Is(err, domain.ErrAssignmentWeeklyHoursInvalid) {
		t.Fatalf("expected ErrAssignmentWeeklyHoursInvalid, got %v", err)
	}
}

func TestCreateAssignmentRepositoryError(t *testing.T) {
	repoErr := errors.New("db error")
	mockRepo := NewMockAssignmentRepository()
	mockRepo.failErr = repoErr
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 8,
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestCreateAssignmentBlocksWhenAssistantHoursExceed22(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      10,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 20,
	})
	if err != nil {
		t.Fatalf("expected no error creating baseline assistant assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      10,
		WorkspaceID: 2,
		Role:        domain.RoleAssistant,
		WeeklyHours: 3,
	})
	if !errors.Is(err, domain.ErrAssignmentAssistantHoursLimitExceeded) {
		t.Fatalf("expected ErrAssignmentAssistantHoursLimitExceeded, got %v", err)
	}
}

func TestCreateAssignmentBlocksWhenMonitorAssignmentsExceed3(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	for i := 1; i <= 3; i++ {
		_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
			UserID:      20,
			WorkspaceID: uint(i),
			Role:        domain.RoleMonitor,
			WeeklyHours: 2,
		})
		if err != nil {
			t.Fatalf("expected no error creating monitor assignment %d, got %v", i, err)
		}
	}

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      20,
		WorkspaceID: 4,
		Role:        domain.RoleMonitor,
		WeeklyHours: 1,
	})
	if !errors.Is(err, domain.ErrAssignmentMonitorCountLimitExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorCountLimitExceeded, got %v", err)
	}
}

func TestCreateAssignmentBlocksWhenMonitorHoursExceed12(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      30,
		WorkspaceID: 1,
		Role:        domain.RoleMonitor,
		WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf("expected no error creating baseline monitor assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      30,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 3,
	})
	if !errors.Is(err, domain.ErrAssignmentMonitorHoursLimitExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorHoursLimitExceeded, got %v", err)
	}
}

func TestCreateAssignmentBlocksWhenMonitorExceedsFortyPercentOfAssistant(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      40,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf("expected no error creating baseline assistant assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      40,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 3,
	})
	if err != nil {
		t.Fatalf("expected no error creating monitor assignment before forty-percent check, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      40,
		WorkspaceID: 3,
		Role:        domain.RoleMonitor,
		WeeklyHours: 2,
	})
	if !errors.Is(err, domain.ErrAssignmentMonitorFortyPercentExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorFortyPercentExceeded, got %v", err)
	}
}

func TestCreateAssignmentAllowsMonitorHoursAtRoundedFortyPercentLimit(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      50,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 11,
	})
	if err != nil {
		t.Fatalf("expected no error creating assistant assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      50,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 4,
	})
	if err != nil {
		t.Fatalf("expected no error creating monitor assignment before rounded limit check, got %v", err)
	}

	output, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      50,
		WorkspaceID: 3,
		Role:        domain.RoleMonitor,
		WeeklyHours: 1,
	})
	if err != nil {
		t.Fatalf("expected no error at rounded forty percent limit, got %v", err)
	}
	if output.ID == 0 {
		t.Fatalf("expected generated id at rounded forty percent limit, got %d", output.ID)
	}
}

func TestCreateAssignmentValidCaseWithMixedRoles(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      60,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 20,
	})
	if err != nil {
		t.Fatalf("expected no error creating assistant assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      60,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 5,
	})
	if err != nil {
		t.Fatalf("expected no error creating first monitor assignment in mixed-role case, got %v", err)
	}

	output, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      60,
		WorkspaceID: 3,
		Role:        domain.RoleMonitor,
		WeeklyHours: 3,
	})
	if err != nil {
		t.Fatalf("expected no error creating valid mixed-role assignment, got %v", err)
	}
	if output.WeeklyHours != 3 {
		t.Fatalf("expected weekly hours 3, got %d", output.WeeklyHours)
	}
}
