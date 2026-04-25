package delivery

import (
	assignmentsDomain "backend/internal/assignments/domain"
	authDomain "backend/internal/auth/domain"
	tasksApplication "backend/internal/tasks/application"
	tasksDelivery "backend/internal/tasks/delivery"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockTaskRepository struct {
	tasks     map[uint]*tasksDomain.Task
	nextID    uint
	createErr error
	updateErr error
}

type mockAssignmentRepository struct {
	assignments map[uint]*assignmentsDomain.Assignment
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks:  make(map[uint]*tasksDomain.Task),
		nextID: 1,
	}
}

func newMockAssignmentRepository() *mockAssignmentRepository {
	return &mockAssignmentRepository{
		assignments: make(map[uint]*assignmentsDomain.Assignment),
	}
}

func (m *mockTaskRepository) Create(task *tasksDomain.Task) error {
	if m.createErr != nil {
		return m.createErr
	}
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
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.tasks[task.ID]; !ok {
		return tasksDomain.ErrTaskNotFound
	}
	m.tasks[task.ID] = task
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

func seedAssignment(repo *mockAssignmentRepository, id, userID uint) {
	repo.assignments[id] = &assignmentsDomain.Assignment{
		ID:          id,
		UserID:      userID,
		WorkspaceID: 7,
		Role:        assignmentsDomain.RoleMonitor,
		WeeklyHours: 4,
	}
}

type mockWorkspaceRepository struct{}

func (m *mockWorkspaceRepository) FindByID(_ uint) (*workspacesDomain.Workspace, error) {
	return &workspacesDomain.Workspace{ID: 7, UserID: 0}, nil
}

func newTaskHandler(repo *mockTaskRepository, assignmentRepo *mockAssignmentRepository, now func() time.Time) *tasksDelivery.TaskHandler {
	return tasksDelivery.NewTaskHandler(
		tasksApplication.NewCreateTask(repo, assignmentRepo, now),
		tasksApplication.NewListTasks(repo),
		tasksApplication.NewGetTaskByID(repo),
		tasksApplication.NewUpdateTask(repo, assignmentRepo, now),
		tasksApplication.NewDeleteTask(repo, now),
		assignmentRepo,
		&mockWorkspaceRepository{},
	)
}

func seedTask(t *testing.T, repo *mockTaskRepository, userID, assignmentID uint, late bool, weekStartDate time.Time) uint {
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
		late,
	)
	if err != nil {
		t.Fatalf("expected seed task, got %v", err)
	}
	if err := repo.Create(task); err != nil {
		t.Fatalf("expected create seed task, got %v", err)
	}
	return task.ID
}

func TestCreateTaskSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	handler := newTaskHandler(taskRepo, assignmentRepo, func() time.Time {
		return time.Date(2026, 4, 9, 9, 0, 0, 0, time.UTC)
	})
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"title":"Prepare class","description":"Review slides","status":"abierto","spent_hours":2,"observations":"","week_start_date":"2026-04-08"}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}
}

func TestCreateTaskBindingError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTaskHandler(newMockTaskRepository(), newMockAssignmentRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"title":"","description":"Review slides","status":"abierto","spent_hours":2,"observations":"","week_start_date":"2026-04-08"}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestCreateTaskReturnsNotFoundForMissingAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTaskHandler(newMockTaskRepository(), newMockAssignmentRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})
	requestBody := bytes.NewBufferString(`{"assignment_id":999,"title":"Prepare class","description":"Review slides","status":"abierto","spent_hours":2,"observations":"","week_start_date":"2026-04-08"}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestGetTaskByIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedTask(t, taskRepo, 2, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	handler := newTaskHandler(taskRepo, assignmentRepo, time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_user", authDomain.AuthenticatedUser{ID: 2, GlobalRole: usersDomain.RoleAdmin})
	request, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.GetTaskByID(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"abierto"`)) {
		t.Fatalf("expected response body to expose canonical status, got %s", recorder.Body.String())
	}
}

func TestUpdateTaskRejectsLateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	seedTask(t, taskRepo, 1, 10, true, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	handler := newTaskHandler(taskRepo, assignmentRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"title":"Prepare class","description":"Review slides","status":"finalizado","spent_hours":3,"observations":"","week_start_date":"2026-04-06"}`)
	request, _ := http.NewRequest(http.MethodPut, "/tasks/1", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.UpdateTask(context)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestCreateTaskRejectsLegacyEnglishStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	handler := newTaskHandler(taskRepo, assignmentRepo, time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"title":"Prepare class","description":"Review slides","status":"open","spent_hours":2,"observations":"","week_start_date":"2026-04-08"}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestCreateTaskInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	taskRepo.createErr = errors.New("db error")
	assignmentRepo := newMockAssignmentRepository()
	seedAssignment(assignmentRepo, 10, 1)
	handler := newTaskHandler(taskRepo, assignmentRepo, time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"title":"Prepare class","description":"Review slides","status":"abierto","spent_hours":2,"observations":"","week_start_date":"2026-04-08"}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}
