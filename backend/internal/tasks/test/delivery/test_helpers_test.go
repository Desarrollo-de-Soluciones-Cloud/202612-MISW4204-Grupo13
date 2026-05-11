package delivery_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	authDomain "backend/internal/auth/domain"
	applicationpkg "backend/internal/tasks/application"
	deliverypkg "backend/internal/tasks/delivery"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"bytes"
	"context"
	"io"
	"testing"
	"time"
)

type mockTaskRepository struct {
	tasks  map[uint]*tasksDomain.Task
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

type mockTaskFileStorage struct {
	uploaded map[string][]byte
	deleted  []string
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{tasks: make(map[uint]*tasksDomain.Task), nextID: 1}
}

func newMockAssignmentRepository() *mockAssignmentRepository {
	return &mockAssignmentRepository{assignments: make(map[uint]*assignmentsDomain.Assignment)}
}

func newMockWorkspaceRepository() *mockWorkspaceRepository {
	return &mockWorkspaceRepository{workspaces: make(map[uint]*workspacesDomain.Workspace)}
}

func newMockWeekRepository() *mockWeekRepository {
	return &mockWeekRepository{weeks: make(map[string]*weeksDomain.Week)}
}

func newMockTaskFileStorage() *mockTaskFileStorage {
	return &mockTaskFileStorage{uploaded: make(map[string][]byte)}
}

func (m *mockTaskRepository) Create(task *tasksDomain.Task) error {
	task.ID = m.nextID
	m.nextID++
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) FindByID(id uint) (*tasksDomain.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return nil, tasksDomain.ErrTaskNotFound
	}
	return task, nil
}

func (m *mockTaskRepository) FindAll() ([]tasksDomain.Task, error) {
	result := make([]tasksDomain.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		result = append(result, *task)
	}
	return result, nil
}

func (m *mockTaskRepository) FindAllByUserID(userID uint) ([]tasksDomain.Task, error) {
	result := make([]tasksDomain.Task, 0)
	for _, task := range m.tasks {
		if task.UserID == userID {
			result = append(result, *task)
		}
	}
	return result, nil
}

func (m *mockTaskRepository) Update(task *tasksDomain.Task) error {
	if _, ok := m.tasks[task.ID]; !ok {
		return tasksDomain.ErrTaskNotFound
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) UpdateAttachments(id uint, attachments []tasksDomain.TaskAttachment) error {
	task, ok := m.tasks[id]
	if !ok {
		return tasksDomain.ErrTaskNotFound
	}
	task.Attachments = attachments
	m.tasks[id] = task
	return nil
}

func (m *mockTaskRepository) Delete(id uint) error {
	if _, ok := m.tasks[id]; !ok {
		return tasksDomain.ErrTaskNotFound
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

func (m *mockTaskFileStorage) Upload(ctx context.Context, objectName string, reader io.Reader, contentType string) error {
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	m.uploaded[objectName] = content
	return nil
}

func (m *mockTaskFileStorage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	content, ok := m.uploaded[objectName]
	if !ok {
		return nil, tasksDomain.ErrTaskNotFound
	}
	return io.NopCloser(bytes.NewReader(content)), nil
}

func (m *mockTaskFileStorage) Delete(ctx context.Context, objectName string) error {
	m.deleted = append(m.deleted, objectName)
	delete(m.uploaded, objectName)
	return nil
}

func seedAssignment(repo *mockAssignmentRepository, id, userID, workspaceID uint) {
	repo.assignments[id] = &assignmentsDomain.Assignment{
		ID:          id,
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        assignmentsDomain.RoleMonitor,
		WeeklyHours: 6,
	}
}

func seedWorkspace(repo *mockWorkspaceRepository, id, userID uint, state workspacesDomain.WorkspaceState) {
	repo.workspaces[id] = &workspacesDomain.Workspace{
		ID:       id,
		UserID:   userID,
		State:    state,
		PeriodID: 1,
	}
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

func seedTask(t *testing.T, repo *mockTaskRepository, userID, assignmentID uint, weekStartDate time.Time, attachments []tasksDomain.TaskAttachment) *tasksDomain.Task {
	t.Helper()
	task, err := tasksDomain.NewTask(
		userID,
		assignmentID,
		nil,
		"Prepare class",
		"Review slides",
		tasksDomain.TaskStatusAbierto,
		2,
		"",
		weekStartDate,
		false,
		attachments,
	)
	if err != nil {
		t.Fatalf("expected seed task, got %v", err)
	}
	if err := repo.Create(task); err != nil {
		t.Fatalf("expected seed task create, got %v", err)
	}
	return task
}

func newTaskHandlerForTest(t *testing.T) (*deliverypkg.TaskHandler, *mockTaskRepository, *mockTaskFileStorage) {
	t.Helper()

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	fileStorage := newMockTaskFileStorage()

	seedAssignment(assignmentRepo, 1, 10, 5)
	seedWorkspace(workspaceRepo, 5, 20, workspacesDomain.ActiveState)
	seedWeek(weekRepo, 1, "2026-04-06")

	now := func() time.Time {
		return time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)
	}

	handler := deliverypkg.NewTaskHandler(
		applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, now),
		applicationpkg.NewListTasks(taskRepo),
		applicationpkg.NewGetTaskByID(taskRepo),
		applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, now),
		applicationpkg.NewSetTaskAttachments(taskRepo),
		applicationpkg.NewDeleteTask(taskRepo, now),
		assignmentRepo,
		workspaceRepo,
		fileStorage,
		"attachments",
	)

	return handler, taskRepo, fileStorage
}

func authenticatedUser(id uint, role usersDomain.UserRole) authDomain.AuthenticatedUser {
	return authDomain.AuthenticatedUser{ID: id, GlobalRole: role}
}
