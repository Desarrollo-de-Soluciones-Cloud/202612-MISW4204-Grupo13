package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	assignmentsInfrastructure "backend/internal/assignments/infrastructure"
	tasksDelivery "backend/internal/tasks/delivery"
	"backend/internal/tasks/application"
	tasksDomain "backend/internal/tasks/domain"
	tasksInfrastructure "backend/internal/tasks/infrastructure"
	authDomain "backend/internal/auth/domain"
	"backend/internal/shared/database"
	usersDomain "backend/internal/users/domain"
	usersInfrastructure "backend/internal/users/infrastructure"
	workspacesDomain "backend/internal/workspaces/domain"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"
	assignmentsDomain "backend/internal/assignments/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTaskHandlerForVisibilityTests(t *testing.T) *tasksDelivery.TaskHandler {
	t.Helper()

	gin.SetMode(gin.TestMode)
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	userRepo := usersInfrastructure.NewUserRepository()
	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()
	assignmentRepo := assignmentsInfrastructure.NewAssignmentRepository()
	taskRepo := tasksInfrastructure.NewTaskRepository()

	if err := userRepo.AutoMigrate(); err != nil {
		t.Fatalf("users automigrate: %v", err)
	}
	if err := workspaceRepo.AutoMigrate(); err != nil {
		t.Fatalf("workspaces automigrate: %v", err)
	}
	if err := assignmentRepo.AutoMigrate(); err != nil {
		t.Fatalf("assignments automigrate: %v", err)
	}
	if err := taskRepo.AutoMigrate(); err != nil {
		t.Fatalf("tasks automigrate: %v", err)
	}

	seedTaskVisibilityData(t, userRepo, workspaceRepo, assignmentRepo, taskRepo)

	createTask := application.NewCreateTask(taskRepo, assignmentRepo, workspaceRepo, nil)
	listTasks := application.NewListTasks(taskRepo)
	getTaskByID := application.NewGetTaskByID(taskRepo)
	updateTask := application.NewUpdateTask(taskRepo, assignmentRepo, nil)
	deleteTask := application.NewDeleteTask(taskRepo, nil)

	return tasksDelivery.NewTaskHandler(
		createTask, listTasks, getTaskByID, updateTask, deleteTask,
		assignmentRepo, workspaceRepo,
	)
}

func seedTaskVisibilityData(
	t *testing.T,
	userRepo *usersInfrastructure.UserRepository,
	workspaceRepo *workspacesInfrastructure.WorkspaceRepository,
	assignmentRepo *assignmentsInfrastructure.AssignmentRepository,
	taskRepo *tasksInfrastructure.TaskRepository,
) {
	t.Helper()

	users := []*usersDomain.User{
		{ID: 1, Name: "Admin", Email: "admin@example.com", Password: "secret123", GlobalRole: usersDomain.RoleAdmin},
		{ID: 10, Name: "Professor One", Email: "prof1@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
		{ID: 11, Name: "Professor Two", Email: "prof2@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
		{ID: 100, Name: "Monitor", Email: "monitor@example.com", Password: "secret123", GlobalRole: usersDomain.RoleMonitor},
		{ID: 200, Name: "Other Monitor", Email: "monitor2@example.com", Password: "secret123", GlobalRole: usersDomain.RoleMonitor},
	}
	for _, u := range users {
		if err := userRepo.Create(u); err != nil {
			t.Fatalf("seed user: %v", err)
		}
	}

	workspaces := []*workspacesDomain.Workspace{
		{ID: 1, PeriodID: 1, UserID: 10, Name: "WS Prof1", Type: workspacesDomain.CourseType, InitialDate: "2026-01-01", FinalDate: "2026-06-30", Observations: "obs", State: workspacesDomain.ActiveState},
		{ID: 2, PeriodID: 1, UserID: 11, Name: "WS Prof2", Type: workspacesDomain.CourseType, InitialDate: "2026-01-01", FinalDate: "2026-06-30", Observations: "obs", State: workspacesDomain.ActiveState},
	}
	for _, ws := range workspaces {
		if err := workspaceRepo.Create(ws); err != nil {
			t.Fatalf("seed workspace: %v", err)
		}
	}

	// Assignment 1: owned by user 100, in workspace 2 (prof 11's)
	// Assignment 2: owned by user 200, in workspace 1 (prof 10's)
	assignments := []*assignmentsDomain.Assignment{
		{ID: 1, UserID: 100, WorkspaceID: 2, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 6},
		{ID: 2, UserID: 200, WorkspaceID: 1, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 6},
	}
	for _, a := range assignments {
		if err := assignmentRepo.Create(a); err != nil {
			t.Fatalf("seed assignment: %v", err)
		}
	}

	// Task 1: owned by user 100, assignment 1 (workspace 2, prof 11)
	// Task 2: owned by user 200, assignment 2 (workspace 1, prof 10)
	tasks := []*tasksDomain.Task{
		{ID: 1, UserID: 100, AssignmentID: 1, Title: "Task One", Description: "desc", Status: tasksDomain.TaskStatusAbierto, SpentHours: 2, Observations: "obs", WeekStartDate: tasksDomain.NormalizeWeekStartDate(time.Now()), Late: false},
		{ID: 2, UserID: 200, AssignmentID: 2, Title: "Task Two", Description: "desc", Status: tasksDomain.TaskStatusAbierto, SpentHours: 2, Observations: "obs", WeekStartDate: time.Date(2026, 12, 7, 0, 0, 0, 0, time.UTC), Late: false},
	}
	for _, task := range tasks {
		if err := database.DB.Create(task).Error; err != nil {
			t.Fatalf("seed task: %v", err)
		}
	}
}

func withTaskCurrentUser(user authDomain.AuthenticatedUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("current_user", user)
		c.Next()
	}
}

// TestMonitorCannotUpdateForeignTask verifies that a monitor cannot update
// a task belonging to a different user.
func TestMonitorCannotUpdateForeignTask(t *testing.T) {
	handler := setupTaskHandlerForVisibilityTests(t)

	router := gin.New()
	// Monitor 100 tries to update task 2 (owned by user 200)
	router.Use(withTaskCurrentUser(authDomain.AuthenticatedUser{ID: 100, GlobalRole: usersDomain.RoleMonitor}))
	router.PUT("/tasks/:id", handler.UpdateTask)

	body := map[string]any{
		"assignment_id":   2,
		"title":           "Updated",
		"description":     "desc",
		"status":          "abierto",
		"spent_hours":     2,
		"observations":    "obs",
		"week_start_date": "2026-12-07",
	}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tasks/2", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d (body: %s)", http.StatusForbidden, w.Code, w.Body.String())
	}
}

// TestProfessorCannotViewTaskFromForeignWorkspace verifies that a professor
// cannot view a task from a workspace they do not own.
func TestProfessorCannotViewTaskFromForeignWorkspace(t *testing.T) {
	handler := setupTaskHandlerForVisibilityTests(t)

	router := gin.New()
	// Professor 10 owns workspace 1; task 1 is in workspace 2 (owned by prof 11)
	router.Use(withTaskCurrentUser(authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor}))
	router.GET("/tasks/:id", handler.GetTaskByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d (body: %s)", http.StatusForbidden, w.Code, w.Body.String())
	}
}

// TestAdminCanDeleteAnyTask verifies that an admin can delete any task.
func TestAdminCanDeleteAnyTask(t *testing.T) {
	handler := setupTaskHandlerForVisibilityTests(t)

	router := gin.New()
	router.Use(withTaskCurrentUser(authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin}))
	router.DELETE("/tasks/:id", handler.DeleteTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d (body: %s)", http.StatusNoContent, w.Code, w.Body.String())
	}
}

// TestMonitorCanUpdateOwnTask verifies that a monitor can update their own task.
func TestMonitorCanUpdateOwnTask(t *testing.T) {
	handler := setupTaskHandlerForVisibilityTests(t)

	router := gin.New()
	// Monitor 100 updates task 1 (owned by user 100, assignment 1)
	router.Use(withTaskCurrentUser(authDomain.AuthenticatedUser{ID: 100, GlobalRole: usersDomain.RoleMonitor}))
	router.PUT("/tasks/:id", handler.UpdateTask)

	body := map[string]any{
		"assignment_id":   1,
		"title":           "Updated Title",
		"description":     "Updated description",
		"status":          "en desarrollo",
		"spent_hours":     3,
		"observations":    "obs",
		"week_start_date": "2026-12-07",
	}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d (body: %s)", http.StatusOK, w.Code, w.Body.String())
	}
}
