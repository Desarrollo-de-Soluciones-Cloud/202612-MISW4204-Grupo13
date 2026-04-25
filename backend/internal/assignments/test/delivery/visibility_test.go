package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	assignmentsInfrastructure "backend/internal/assignments/infrastructure"
	"backend/internal/assignments/application"
	assignmentsDelivery "backend/internal/assignments/delivery"
	assignmentsDomain "backend/internal/assignments/domain"
	authDomain "backend/internal/auth/domain"
	"backend/internal/shared/database"
	usersDomain "backend/internal/users/domain"
	usersInfrastructure "backend/internal/users/infrastructure"
	workspacesDomain "backend/internal/workspaces/domain"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAssignmentHandlerForVisibilityTests(t *testing.T) *assignmentsDelivery.AssignmentHandler {
	t.Helper()

	gin.SetMode(gin.TestMode)
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	assignmentRepo := assignmentsInfrastructure.NewAssignmentRepository()
	userRepo := usersInfrastructure.NewUserRepository()
	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()

	if err := userRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected users automigrate, got %v", err)
	}
	if err := workspaceRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected workspaces automigrate, got %v", err)
	}
	if err := assignmentRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected assignments automigrate, got %v", err)
	}

	seedVisibilityData(t, userRepo, workspaceRepo, assignmentRepo)

	createAssignment := application.NewCreateAssignment(assignmentRepo).WithRepositories(userRepo, workspaceRepo)
	getAssignmentByID := application.NewGetAssignmentByID(assignmentRepo)
	listAssignmentsByUserID := application.NewListAssignmentsByUserID(assignmentRepo)
	updateAssignment := application.NewUpdateAssignment(assignmentRepo)

	return assignmentsDelivery.NewAssignmentHandler(
		createAssignment,
		getAssignmentByID,
		listAssignmentsByUserID,
		updateAssignment,
		workspaceRepo,
	)
}

func seedVisibilityData(t *testing.T, userRepo *usersInfrastructure.UserRepository, workspaceRepo *workspacesInfrastructure.WorkspaceRepository, assignmentRepo *assignmentsInfrastructure.AssignmentRepository) {
	t.Helper()

	users := []*usersDomain.User{
		{ID: 1, Name: "Admin", Email: "admin@example.com", Password: "secret123", GlobalRole: usersDomain.RoleAdmin},
		{ID: 10, Name: "Professor One", Email: "prof1@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
		{ID: 11, Name: "Professor Two", Email: "prof2@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
		{ID: 100, Name: "Monitor", Email: "monitor@example.com", Password: "secret123", GlobalRole: usersDomain.RoleMonitor},
		{ID: 101, Name: "Assistant", Email: "assistant@example.com", Password: "secret123", GlobalRole: usersDomain.RoleAssistant},
		{ID: 200, Name: "Student A", Email: "studenta@example.com", Password: "secret123", GlobalRole: usersDomain.RoleAssistant},
		{ID: 201, Name: "Student B", Email: "studentb@example.com", Password: "secret123", GlobalRole: usersDomain.RoleMonitor},
	}

	for _, user := range users {
		if err := userRepo.Create(user); err != nil {
			t.Fatalf("expected user seed, got %v", err)
		}
	}

	workspaceOwnedByProf1 := &workspacesDomain.Workspace{
		ID:           1,
		PeriodID:     1,
		UserID:       10,
		Name:         "Workspace Prof 1",
		Type:         workspacesDomain.CourseType,
		InitialDate:  "2026-01-01",
		FinalDate:    "2026-06-30",
		Observations: "obs",
		State:        workspacesDomain.ActiveState,
	}
	workspaceOwnedByProf2 := &workspacesDomain.Workspace{
		ID:           2,
		PeriodID:     1,
		UserID:       11,
		Name:         "Workspace Prof 2",
		Type:         workspacesDomain.CourseType,
		InitialDate:  "2026-01-01",
		FinalDate:    "2026-06-30",
		Observations: "obs",
		State:        workspacesDomain.ActiveState,
	}

	if err := workspaceRepo.Create(workspaceOwnedByProf1); err != nil {
		t.Fatalf("expected workspace seed, got %v", err)
	}
	if err := workspaceRepo.Create(workspaceOwnedByProf2); err != nil {
		t.Fatalf("expected workspace seed, got %v", err)
	}

	assignmentInProf2Workspace, err := assignmentsDomain.NewAssignment(201, 2, assignmentsDomain.RoleMonitor, 6)
	if err != nil {
		t.Fatalf("expected assignment seed build, got %v", err)
	}
	if err := assignmentRepo.Create(assignmentInProf2Workspace); err != nil {
		t.Fatalf("expected assignment seed create, got %v", err)
	}
}

func withCurrentUser(user authDomain.AuthenticatedUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("current_user", user)
		c.Next()
	}
}

func TestMonitorCannotViewAssignmentFromAnotherUser(t *testing.T) {
	handler := setupAssignmentHandlerForVisibilityTests(t)

	router := gin.New()
	router.Use(withCurrentUser(authDomain.AuthenticatedUser{ID: 100, GlobalRole: usersDomain.RoleMonitor}))
	router.GET("/assignments/:id", handler.GetAssignmentByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/assignments/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestProfessorCannotViewAssignmentFromForeignWorkspace(t *testing.T) {
	handler := setupAssignmentHandlerForVisibilityTests(t)

	router := gin.New()
	router.Use(withCurrentUser(authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor}))
	router.GET("/assignments/:id", handler.GetAssignmentByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/assignments/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestAdminCanViewAnyAssignment(t *testing.T) {
	handler := setupAssignmentHandlerForVisibilityTests(t)

	router := gin.New()
	router.Use(withCurrentUser(authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin}))
	router.GET("/assignments/:id", handler.GetAssignmentByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/assignments/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProfessorCanCreateAssignmentOnlyInOwnWorkspace(t *testing.T) {
	handler := setupAssignmentHandlerForVisibilityTests(t)

	router := gin.New()
	router.Use(withCurrentUser(authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor}))
	router.POST("/assignments", handler.CreateAssignment)

	foreignWorkspaceBody := map[string]any{
		"user_id":      200,
		"workspace_id": 2,
		"role":         string(assignmentsDomain.RoleAssistant),
		"weekly_hours": 8,
	}
	foreignJSON, _ := json.Marshal(foreignWorkspaceBody)

	forbiddenWriter := httptest.NewRecorder()
	forbiddenRequest, _ := http.NewRequest(http.MethodPost, "/assignments", bytes.NewBuffer(foreignJSON))
	forbiddenRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(forbiddenWriter, forbiddenRequest)

	if forbiddenWriter.Code != http.StatusForbidden {
		t.Fatalf("expected status %d for foreign workspace, got %d", http.StatusForbidden, forbiddenWriter.Code)
	}

	ownWorkspaceBody := map[string]any{
		"user_id":      200,
		"workspace_id": 1,
		"role":         string(assignmentsDomain.RoleAssistant),
		"weekly_hours": 8,
	}
	ownJSON, _ := json.Marshal(ownWorkspaceBody)

	createdWriter := httptest.NewRecorder()
	createdRequest, _ := http.NewRequest(http.MethodPost, "/assignments", bytes.NewBuffer(ownJSON))
	createdRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(createdWriter, createdRequest)

	if createdWriter.Code != http.StatusCreated {
		t.Fatalf("expected status %d for own workspace, got %d", http.StatusCreated, createdWriter.Code)
	}
}
