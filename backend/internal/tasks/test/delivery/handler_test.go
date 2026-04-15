package delivery

import (
	assignmentsDomain "backend/internal/assignments/domain"
	tasksApplication "backend/internal/tasks/application"
	tasksDelivery "backend/internal/tasks/delivery"
	tasksDomain "backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
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

type mockWorkspaceRepository struct {
	workspaces map[uint]*workspacesDomain.Workspace
}

type mockWeekRepository struct {
	weeks map[uint]*weeksDomain.Week
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks:  make(map[uint]*tasksDomain.Task),
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

func (m *mockTaskRepository) Create(task *tasksDomain.Task) error {
	if m.createErr != nil {
		return m.createErr
	}
	task.ID = m.nextID
	for i := range task.Attachments {
		task.Attachments[i].ID = uint(i + 1)
		task.Attachments[i].TaskID = task.ID
	}
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
	for i := range task.Attachments {
		task.Attachments[i].ID = uint(i + 1)
		task.Attachments[i].TaskID = task.ID
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
		WeeklyHours: 4,
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

func newTaskHandler(
	repo *mockTaskRepository,
	assignmentRepo *mockAssignmentRepository,
	workspaceRepo *mockWorkspaceRepository,
	weekRepo *mockWeekRepository,
	now func() time.Time,
) *tasksDelivery.TaskHandler {
	return tasksDelivery.NewTaskHandler(
		tasksApplication.NewCreateTask(repo, assignmentRepo, workspaceRepo, weekRepo, now),
		tasksApplication.NewListTasks(repo),
		tasksApplication.NewGetTaskByID(repo),
		tasksApplication.NewUpdateTask(repo, assignmentRepo, workspaceRepo, weekRepo, now),
		tasksApplication.NewDeleteTask(repo, weekRepo, now),
	)
}

func seedTask(t *testing.T, repo *mockTaskRepository, userID, assignmentID, weekID uint, late bool, weekStartDate time.Time) uint {
	t.Helper()

	task, err := tasksDomain.NewTask(
		userID,
		assignmentID,
		weekID,
		"Prepare class",
		"Review slides",
		tasksDomain.TaskStatusAbierto,
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
	return task.ID
}

func TestCreateTaskSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 7)
	seedWorkspace(workspaceRepo, 7, 9)
	seedWeek(weekRepo, 20, 9, "2026-04-13", "2026-04-19")
	handler := newTaskHandler(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"week_id":20,"title":"Prepare class","description":"Review slides","status":"abierto","spent_hours":2,"observations":"","attachments":[{"path":"docs/guide.pdf"}]}`)
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

	handler := newTaskHandler(newMockTaskRepository(), newMockAssignmentRepository(), newMockWorkspaceRepository(), newMockWeekRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"title":"","description":"Review slides","status":"abierto","spent_hours":2}`)
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

	handler := newTaskHandler(newMockTaskRepository(), newMockAssignmentRepository(), newMockWorkspaceRepository(), newMockWeekRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"assignment_id":999,"week_id":20,"title":"Prepare class","description":"Review slides","status":"abierto","spent_hours":2}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestCreateTaskRejectsInvalidAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 7)
	seedWorkspace(workspaceRepo, 7, 9)
	seedWeek(weekRepo, 20, 9, "2026-04-13", "2026-04-19")
	handler := newTaskHandler(taskRepo, assignmentRepo, workspaceRepo, weekRepo, time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"week_id":20,"title":"Prepare class","description":"Review slides","status":"abierto","spent_hours":2,"attachments":[{"path":" "} ]}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestGetTaskByIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	seedTask(t, taskRepo, 2, 20, 30, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))
	handler := newTaskHandler(taskRepo, newMockAssignmentRepository(), newMockWorkspaceRepository(), newMockWeekRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
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

func TestGetTaskByIDBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTaskHandler(newMockTaskRepository(), newMockAssignmentRepository(), newMockWorkspaceRepository(), newMockWeekRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/tasks/abc", nil)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler.GetTaskByID(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestListTasksSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	seedTask(t, taskRepo, 2, 20, 30, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))
	handler := newTaskHandler(taskRepo, newMockAssignmentRepository(), newMockWorkspaceRepository(), newMockWeekRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	context.Request = request

	handler.ListTasks(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestUpdateTaskRejectsClosedWeek(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 7)
	seedWorkspace(workspaceRepo, 7, 9)
	seedWeek(weekRepo, 20, 9, "2026-04-06", "2026-04-12")
	seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC))
	handler := newTaskHandler(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"week_id":20,"title":"Prepare class","description":"Review slides","status":"finalizado","spent_hours":3}`)
	request, _ := http.NewRequest(http.MethodPut, "/tasks/1", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.UpdateTask(context)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestUpdateTaskRejectsChangingAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 7)
	seedAssignment(assignmentRepo, 11, 1, 7)
	seedWorkspace(workspaceRepo, 7, 9)
	seedWeek(weekRepo, 20, 9, "2026-04-13", "2026-04-19")
	seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))
	handler := newTaskHandler(taskRepo, assignmentRepo, workspaceRepo, weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"assignment_id":11,"week_id":20,"title":"Prepare class","description":"Review slides","status":"finalizado","spent_hours":3}`)
	request, _ := http.NewRequest(http.MethodPut, "/tasks/1", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.UpdateTask(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestDeleteTaskSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	weekRepo := newMockWeekRepository()
	seedWeek(weekRepo, 20, 9, "2026-04-13", "2026-04-19")
	seedTask(t, taskRepo, 1, 10, 20, false, time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC))
	handler := newTaskHandler(taskRepo, newMockAssignmentRepository(), newMockWorkspaceRepository(), weekRepo, func() time.Time {
		return time.Date(2026, 4, 14, 9, 0, 0, 0, time.UTC)
	})
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.DeleteTask(context)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", recorder.Code)
	}
}

func TestDeleteTaskBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newTaskHandler(newMockTaskRepository(), newMockAssignmentRepository(), newMockWorkspaceRepository(), newMockWeekRepository(), time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler.DeleteTask(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestCreateTaskInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := newMockTaskRepository()
	taskRepo.createErr = errors.New("db error")
	assignmentRepo := newMockAssignmentRepository()
	workspaceRepo := newMockWorkspaceRepository()
	weekRepo := newMockWeekRepository()
	seedAssignment(assignmentRepo, 10, 1, 7)
	seedWorkspace(workspaceRepo, 7, 9)
	seedWeek(weekRepo, 20, 9, "2026-04-13", "2026-04-19")
	handler := newTaskHandler(taskRepo, assignmentRepo, workspaceRepo, weekRepo, time.Now)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"assignment_id":10,"week_id":20,"title":"Prepare class","description":"Review slides","status":"abierto","spent_hours":2}`)
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateTask(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}
