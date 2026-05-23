package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	"backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"testing"
	"time"
)

type mockTaskRepository struct {
	tasks  map[uint]*domain.Task
	nextID uint
}

type mockAssignmentRepository struct {
	assignments map[uint]*assignmentsDomain.Assignment
}

type mockWorkspaceRepository struct {
	workspaces map[uint]*workspacesDomain.Workspace
}

type mockWeekRepository struct {
	weeks map[string]*weeksDomain.Week
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks:  make(map[uint]*domain.Task),
		nextID: 1,
	}
}

func newMockAssignmentRepository() *mockAssignmentRepository {
	return &mockAssignmentRepository{
		assignments: make(map[uint]*assignmentsDomain.Assignment),
	}
}

func newMockWorkspaceRepository() *mockWorkspaceRepository {
	return &mockWorkspaceRepository{workspaces: make(map[uint]*workspacesDomain.Workspace)}
}

func newMockWeekRepository() *mockWeekRepository {
	return &mockWeekRepository{weeks: make(map[string]*weeksDomain.Week)}
}

func (m *mockTaskRepository) Create(task *domain.Task) error {
	task.ID = m.nextID
	m.nextID++
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) FindByID(id uint) (*domain.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	return task, nil
}

func (m *mockTaskRepository) FindAll() ([]domain.Task, error) {
	result := make([]domain.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		result = append(result, *task)
	}
	return result, nil
}

func (m *mockTaskRepository) FindAllByUserID(userID uint) ([]domain.Task, error) {
	result := make([]domain.Task, 0)
	for _, task := range m.tasks {
		if task.UserID == userID {
			result = append(result, *task)
		}
	}
	return result, nil
}

func (m *mockTaskRepository) Update(task *domain.Task) error {
	if _, ok := m.tasks[task.ID]; !ok {
		return domain.ErrTaskNotFound
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) UpdateAttachments(id uint, attachments []domain.TaskAttachment) error {
	task, ok := m.tasks[id]
	if !ok {
		return domain.ErrTaskNotFound
	}
	task.Attachments = attachments
	m.tasks[id] = task
	return nil
}

func (m *mockTaskRepository) Delete(id uint) error {
	if _, ok := m.tasks[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(m.tasks, id)
	return nil
}

func (m *mockAssignmentRepository) FindByID(id uint) (*assignmentsDomain.Assignment, error) {
	assignment, ok := m.assignments[id]
	if !ok {
		return nil, assignmentsDomain.ErrAssignmentNotFound
	}
	return assignment, nil
}

func (m *mockWorkspaceRepository) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	workspace, ok := m.workspaces[id]
	if !ok {
		return nil, workspacesDomain.ErrWorkspaceNotFound
	}
	return workspace, nil
}

func (m *mockWeekRepository) FindByPeriodIDAndStartDate(periodID uint, startDate string) (*weeksDomain.Week, error) {
	week, ok := m.weeks[startDate]
	if !ok {
		return nil, weeksDomain.ErrWeekNotFound
	}
	return week, nil
}

func seedAssignment(repo *mockAssignmentRepository, id, userID uint) {
	repo.assignments[id] = &assignmentsDomain.Assignment{
		ID:          id,
		UserID:      userID,
		WorkspaceID: 5,
		Role:        assignmentsDomain.RoleMonitor,
		WeeklyHours: 6,
	}
}

func seedWorkspace(repo *mockWorkspaceRepository, id uint, state workspacesDomain.WorkspaceState) {
	repo.workspaces[id] = &workspacesDomain.Workspace{ID: id, State: state, PeriodID: 1}
}

func seedWeek(repo *mockWeekRepository, periodID uint, initialDate string) {
	repo.weeks[initialDate] = &weeksDomain.Week{
		ID:          1,
		PeriodID:    periodID,
		Number:      1,
		InitialDate: initialDate,
		FinalDate:   "2026-04-12",
	}
}

func seedTask(t *testing.T, repo *mockTaskRepository, userID, assignmentID uint, late bool, weekStartDate time.Time) *domain.Task {
	t.Helper()

	task, err := domain.NewTask(domain.TaskInput{
		UserID:        userID,
		AssignmentID:  assignmentID,
		WeekID:        nil,
		Title:         "Prepare class",
		Description:   "Review slides",
		Status:        domain.TaskStatusAbierto,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: weekStartDate,
		Late:          late,
		Attachments:   []domain.TaskAttachment{},
	})
	if err != nil {
		t.Fatalf("expected seed task, got %v", err)
	}
	if err := repo.Create(task); err != nil {
		t.Fatalf("expected create seed task, got %v", err)
	}
	return task
}
