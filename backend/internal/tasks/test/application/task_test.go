package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	applicationpkg "backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
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

type mockWorkspaceRepository struct {
	workspaces map[uint]*workspacesDomain.Workspace
}

type mockWeekRepository struct {
	weeks map[uint]*weeksDomain.Week
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks:  make(map[uint]*domain.Task),
		nextID: 1,
	}
}

func newMockAssignmentRepository() *mockAssignmentRepository {
	return &mockAssignmentRepository{assignments: make(map[uint]*assignmentsDomain.Assignment)}
}

func newMockWorkspaceRepository() *mockWorkspaceRepository {
	return &mockWorkspaceRepository{workspaces: make(map[uint]*workspacesDomain.Workspace)}
}

func newMockWeekRepository() *mockWeekRepository {
	return &mockWeekRepository{weeks: make(map[uint]*weeksDomain.Week)}
}

func (m *mockTaskRepository) Create(task *domain.Task) error {
	task.ID = m.nextID
	for i := range task.Attachments {
		task.Attachments[i].ID = uint(i + 1)
		task.Attachments[i].TaskID = task.ID
	}
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
	for i := range task.Attachments {
		task.Attachments[i].ID = uint(i + 1)
		task.Attachments[i].TaskID = task.ID
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

func (m *mockWorkspaceRepository) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	workspace, ok := m.workspaces[id]
	if !ok {
		return nil, workspacesDomain.ErrWorkspaceNotFound
	}
	return workspace, nil
}

func (m *mockWeekRepository) FindByID(id uint) (*weeksDomain.Week, error) {
	week, ok := m.weeks[id]
	if !ok {
		return nil, weeksDomain.ErrWeekNotFound
	}
	return week, nil
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

func seedWorkspace(repo *mockWorkspaceRepository, id, periodID uint) {
	repo.workspaces[id] = &workspacesDomain.Workspace{
		ID:       id,
		PeriodID: periodID,
	}
}

func seedWeek(repo *mockWeekRepository, id, periodID uint, initialDate, finalDate string) {
	repo.weeks[id] = &weeksDomain.Week{
		ID:          id,
		PeriodID:    periodID,
		Number:      1,
		InitialDate: initialDate,
		FinalDate:   finalDate,
	}
}

func seedTask(t *testing.T, repo *mockTaskRepository, userID, assignmentID, weekID uint, late bool, weekStartDate time.Time) *domain.Task {
	t.Helper()

	task, err := domain.NewTask(
		userID,
		assignmentID,
		weekID,
		"Prepare class",
		"Review slides",
		domain.TaskStatusAbierto,
		2,
		"",
		nil,
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

func TestCreateTaskBuildsTaskFromAssignmentAndWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 3, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 8, "2026-04-13", "2026-04-19")
	now := func() time.Time { return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC) }

	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, now)
	output, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID: 10,
		WeekID:       20,
		Title:        "Prepare class",
		Description:  "Review slides",
		Status:       domain.TaskStatusAbierto,
		SpentHours:   2,
		Observations: "",
		Attachments:  []applicationpkg.TaskAttachmentInput{{Path: "docs/plan.pdf"}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.AssignmentID != 10 || output.WeekID != 20 {
		t.Fatalf("expected assignment 10 and week 20, got %d and %d", output.AssignmentID, output.WeekID)
	}
	if output.UserID != 3 {
		t.Fatalf("expected user id 3, got %d", output.UserID)
	}
	if output.Late {
		t.Fatal("expected task not to be late")
	}
	if output.WeekStartDate.Format("2006-01-02") != "2026-04-13" {
		t.Fatalf("expected week start 2026-04-13, got %s", output.WeekStartDate.Format("2006-01-02"))
	}
	if len(output.Attachments) != 1 || output.Attachments[0].Path != "docs/plan.pdf" {
		t.Fatalf("expected attachments to be preserved, got %+v", output.Attachments)
	}
}

func TestCreateTaskMarksClosedWeekAsLate(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 3, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 8, "2026-04-06", "2026-04-12")

	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})

	output, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID: 10,
		WeekID:       20,
		Title:        "Prepare class",
		Description:  "Review slides",
		Status:       domain.TaskStatusAbierto,
		SpentHours:   2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !output.Late {
		t.Fatal("expected task to be marked as late")
	}
}

func TestCreateTaskRejectsMissingAssignment(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, time.Now)

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID: 10,
		WeekID:       20,
		Title:        "Prepare class",
		Description:  "Review slides",
		Status:       domain.TaskStatusAbierto,
		SpentHours:   2,
	})
	if !errors.Is(err, domain.ErrTaskAssignmentNotFound) {
		t.Fatalf("expected ErrTaskAssignmentNotFound, got %v", err)
	}
}

func TestCreateTaskRejectsLegacyEnglishStatus(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 3, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 8, "2026-04-13", "2026-04-19")
	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, time.Now)

	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID: 10,
		WeekID:       20,
		Title:        "Prepare class",
		Description:  "Review slides",
		Status:       domain.TaskStatus("open"),
		SpentHours:   2,
	})
	if !errors.Is(err, domain.ErrTaskStatusInvalid) {
		t.Fatalf("expected ErrTaskStatusInvalid, got %v", err)
	}
}

func TestCreateTaskRejectsWeekFromDifferentPeriod(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 3, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 99, "2026-04-13", "2026-04-19")

	createTask := applicationpkg.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, time.Now)
	_, err := createTask.Execute(applicationpkg.CreateTaskInput{
		AssignmentID: 10,
		WeekID:       20,
		Title:        "Prepare class",
		Description:  "Review slides",
		Status:       domain.TaskStatusAbierto,
		SpentHours:   2,
	})
	if !errors.Is(err, domain.ErrTaskWeekPeriodMismatch) {
		t.Fatalf("expected ErrTaskWeekPeriodMismatch, got %v", err)
	}
}

func TestListTasksReturnsAllTasks(t *testing.T) {
	taskRepo := newMockTaskRepository()
	seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))
	seedTask(t, taskRepo, 2, 20, 30, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))

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
	task := seedTask(t, taskRepo, 2, 20, 30, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))
	getTask := applicationpkg.NewGetTaskByID(taskRepo)

	output, err := getTask.Execute(applicationpkg.GetTaskByIDInput{ID: task.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.ID != task.ID {
		t.Fatalf("expected task id %d, got %d", task.ID, output.ID)
	}
}

func TestGetTaskByIDReturnsNotFound(t *testing.T) {
	taskRepo := newMockTaskRepository()
	getTask := applicationpkg.NewGetTaskByID(taskRepo)

	_, err := getTask.Execute(applicationpkg.GetTaskByIDInput{ID: 999})
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestUpdateTaskRejectsChangingWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 8, "2026-04-13", "2026-04-19")
	seedWeek(weekRepo, 21, 8, "2026-04-20", "2026-04-26")
	task := seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))

	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	_, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:           task.ID,
		AssignmentID: 10,
		WeekID:       21,
		Title:        "Prepare class 2",
		Description:  "Review slides",
		Status:       domain.TaskStatusFinalizado,
		SpentHours:   3,
	})
	if !errors.Is(err, domain.ErrTaskWeekChangeForbidden) {
		t.Fatalf("expected ErrTaskWeekChangeForbidden, got %v", err)
	}
}

func TestUpdateTaskRejectsChangingAssignment(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 5)
	seedAssignment(assignmentRepo, 11, 1, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 8, "2026-04-13", "2026-04-19")
	task := seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))

	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	_, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:           task.ID,
		AssignmentID: 11,
		WeekID:       20,
		Title:        "Prepare class 2",
		Description:  "Review slides",
		Status:       domain.TaskStatusFinalizado,
		SpentHours:   3,
	})
	if !errors.Is(err, domain.ErrTaskAssignmentChangeForbidden) {
		t.Fatalf("expected ErrTaskAssignmentChangeForbidden, got %v", err)
	}
}

func TestUpdateTaskRejectsInactiveWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 8, "2026-04-06", "2026-04-12")
	task := seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))

	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	_, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:           task.ID,
		AssignmentID: 10,
		WeekID:       20,
		Title:        "Prepare class 2",
		Description:  "Review slides",
		Status:       domain.TaskStatusFinalizado,
		SpentHours:   3,
	})
	if !errors.Is(err, domain.ErrTaskLateUpdateForbidden) {
		t.Fatalf("expected ErrTaskLateUpdateForbidden, got %v", err)
	}
}

func TestUpdateTaskSuccess(t *testing.T) {
	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 5)
	seedWorkspace(workspaceRepo, 5, 8)
	seedWeek(weekRepo, 20, 8, "2026-04-13", "2026-04-19")
	task := seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))

	updateTask := applicationpkg.NewUpdateTask(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	output, err := updateTask.Execute(applicationpkg.UpdateTaskInput{
		ID:           task.ID,
		AssignmentID: 10,
		WeekID:       20,
		Title:        "Prepare class 2",
		Description:  "Review slides updated",
		Status:       domain.TaskStatusFinalizado,
		SpentHours:   3,
		Attachments:  []applicationpkg.TaskAttachmentInput{{Path: "docs/result.pdf"}},
	})
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}
	if output.Title != "Prepare class 2" || len(output.Attachments) != 1 {
		t.Fatalf("expected updated task output, got %+v", output)
	}
}

func TestDeleteTaskRejectsInactiveWeek(t *testing.T) {
	taskRepo := newMockTaskRepository()
	weekRepo := newMockWeekRepository()
	seedWeek(weekRepo, 20, 8, "2026-04-06", "2026-04-12")
	task := seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))

	deleteTask := applicationpkg.NewDeleteTask(taskRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	err := deleteTask.Execute(applicationpkg.DeleteTaskInput{ID: task.ID})
	if !errors.Is(err, domain.ErrTaskDeleteForbidden) {
		t.Fatalf("expected ErrTaskDeleteForbidden, got %v", err)
	}
}

func TestDeleteTaskSuccess(t *testing.T) {
	taskRepo := newMockTaskRepository()
	weekRepo := newMockWeekRepository()
	seedWeek(weekRepo, 20, 8, "2026-04-13", "2026-04-19")
	task := seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))

	deleteTask := applicationpkg.NewDeleteTask(taskRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	if err := deleteTask.Execute(applicationpkg.DeleteTaskInput{ID: task.ID}); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}
	if _, err := taskRepo.FindByID(task.ID); !errors.Is(err, domain.ErrTaskNotFound) {
		t.Fatalf("expected task to be deleted, got %v", err)
	}
}
