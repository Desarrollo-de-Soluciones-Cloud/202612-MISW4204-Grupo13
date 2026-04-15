package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	applicationpkg "backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	"errors"
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

func seedAssignment(repo *mockAssignmentRepository, id, userID uint) {
	repo.assignments[id] = &assignmentsDomain.Assignment{
		ID:          id,
		UserID:      userID,
		WorkspaceID: 5,
		Role:        assignmentsDomain.RoleMonitor,
		WeeklyHours: 6,
	}
}

func seedTask(t *testing.T, repo *mockTaskRepository, userID, assignmentID uint, late bool, weekStartDate time.Time) *domain.Task {
	t.Helper()

	task, err := domain.NewTask(
		userID,
		assignmentID,
		nil,
		"Prepare class",
		"Review slides",
		domain.TaskStatusAbierto,
		2,
		"",
		weekStartDate,
		late,
	)
	if err != nil {
		t.Fatalf("expected seed task, got %v", err)
	}
	if err := repo.Create(task); err != nil {
		t.Fatalf("expected create seed task, got %v", err)
	}
	return task
}

func TestCreateTaskMarksPastWeekAsLate(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 3)
	now := func() time.Time { return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC) }
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, now)

	output, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID:  10,
		Title:         "Prepare class",
		Description:   "Review slides",
		Status:        domain.TaskStatusAbierto,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !output.Late {
		t.Fatal("expected task to be marked as late")
	}
	if output.AssignmentID != 10 {
		t.Fatalf("expected assignment id 10, got %d", output.AssignmentID)
	}
}

func TestCreateTaskRejectsMissingAssignment(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, time.Now)

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID:  10,
		Title:         "Prepare class",
		Description:   "Review slides",
		Status:        domain.TaskStatusAbierto,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrTaskAssignmentNotFound) {
		t.Fatalf("expected ErrTaskAssignmentNotFound, got %v", err)
	}
}

func TestCreateTaskRejectsLegacyEnglishStatus(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 3)
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, time.Now)

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID:  10,
		Title:         "Prepare class",
		Description:   "Review slides",
		Status:        domain.TaskStatus("open"),
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrTaskStatusInvalid) {
		t.Fatalf("expected ErrTaskStatusInvalid, got %v", err)
	}
}

func TestListTasksReturnsAllTasks(t *testing.T) {
	taskRepo := newMockTaskRepository()
	seedTask(t, taskRepo, 1, 10, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	seedTask(t, taskRepo, 2, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))

	listTasks := applicationpkg.NewListTasks(taskRepo)
	output, err := listTasks.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(output.Tasks))
	}
}

func TestGetTaskByIDReturnsTask(t *testing.T) {
	taskRepo := newMockTaskRepository()
	task := seedTask(t, taskRepo, 2, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	getTask := applicationpkg.NewGetTaskByID(taskRepo)

	output, err := getTask.Execute(applicationpkg.GetTaskByIDInput{ID: task.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.ID != task.ID {
		t.Fatalf("expected task id %d, got %d", task.ID, output.ID)
	}
}

func TestUpdateTaskRejectsLateTask(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	task := seedTask(t, taskRepo, 1, 10, true, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})

	_, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:            task.ID,
		AssignmentID:  10,
		Title:         "Prepare class 2",
		Description:   "Review slides",
		Status:        domain.TaskStatusFinalizado,
		SpentHours:    3,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domain.ErrTaskLateUpdateForbidden) {
		t.Fatalf("expected ErrTaskLateUpdateForbidden, got %v", err)
	}
}

func TestUpdateTaskAllowsChangingStatusDuringActiveWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	task := seedTask(t, taskRepo, 1, 10, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, func() time.Time {
		return time.Date(2026, 4, 9, 9, 0, 0, 0, time.UTC)
	})

	output, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:            task.ID,
		AssignmentID:  10,
		Title:         "Prepare class",
		Description:   "Review slides",
		Status:        domain.TaskStatusFinalizado,
		SpentHours:    2,
		Observations:  "",
		WeekStartDate: time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Status != domain.TaskStatusFinalizado {
		t.Fatalf("expected updated status %q, got %q", domain.TaskStatusFinalizado, output.Status)
	}
}

func TestDeleteTaskRejectsInactiveWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	task := seedTask(t, taskRepo, 1, 10, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	deleteTask := applicationpkg.NewDeleteTask(taskRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})

	err := deleteTask.Execute(applicationpkg.DeleteTaskInput{ID: task.ID})
	if !errors.Is(err, domain.ErrTaskDeleteForbidden) {
		t.Fatalf("expected ErrTaskDeleteForbidden, got %v", err)
	}
}
