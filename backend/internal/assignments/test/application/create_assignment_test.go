package application

import (
	applicationpkg "backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	"errors"
	"testing"
)

type MockAssignmentRepository struct {
	byID    map[uint]*domain.Assignment
	nextID  uint
	failErr error
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

func (m *MockAssignmentRepository) Update(assignment *domain.Assignment) error {
	if _, ok := m.byID[assignment.ID]; !ok {
		return domain.ErrAssignmentNotFound
	}
	m.byID[assignment.ID] = assignment
	return nil
}

func TestCreateAssignmentSuccess(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)

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
