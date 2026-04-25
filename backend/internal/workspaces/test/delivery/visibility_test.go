package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	authDomain "backend/internal/auth/domain"
	periodsDomain "backend/internal/periods/domain"
	periodsInfrastructure "backend/internal/periods/infrastructure"
	"backend/internal/shared/database"
	usersDomain "backend/internal/users/domain"
	usersInfrastructure "backend/internal/users/infrastructure"
	workspacesApplication "backend/internal/workspaces/application"
	workspacesDelivery "backend/internal/workspaces/delivery"
	workspacesDomain "backend/internal/workspaces/domain"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWorkspaceHandlerForVisibilityTests(t *testing.T) *workspacesDelivery.WorkspaceHandler {
	t.Helper()

	gin.SetMode(gin.TestMode)
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	periodRepo := periodsInfrastructure.NewPeriodRepository()
	userRepo := usersInfrastructure.NewUserRepository()
	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()

	if err := periodRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected periods automigrate, got %v", err)
	}
	if err := userRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected users automigrate, got %v", err)
	}
	if err := workspaceRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected workspaces automigrate, got %v", err)
	}

	seedWorkspaceVisibilityData(t, periodRepo, userRepo, workspaceRepo)

	createWorkspace := workspacesApplication.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	listWorkspaces := workspacesApplication.NewListWorkspaces(workspaceRepo)
	listWorkspacesByPeriod := workspacesApplication.NewListWorkspacesByPeriod(workspaceRepo)
	getWorkspaceByID := workspacesApplication.NewGetWorkspaceByID(workspaceRepo)
	updateWorkspace := workspacesApplication.NewUpdateWorkspace(workspaceRepo, periodRepo, userRepo)
	deleteWorkspace := workspacesApplication.NewDeleteWorkspace(workspaceRepo)

	return workspacesDelivery.NewWorkspaceHandler(createWorkspace, listWorkspaces, listWorkspacesByPeriod, getWorkspaceByID, updateWorkspace, deleteWorkspace)
}

func seedWorkspaceVisibilityData(t *testing.T, periodRepo *periodsInfrastructure.PeriodRepository, userRepo *usersInfrastructure.UserRepository, workspaceRepo *workspacesInfrastructure.WorkspaceRepository) {
	t.Helper()

	period := &periodsDomain.Period{ID: 1, Name: "2027-10", InitialDate: "2027-01-05", FinalDate: "2027-04-26", InscriptionFinalDate: time.Now().AddDate(0, 0, 10).Format("2006-01-02"), WeeksCount: 16, PeriodState: periodsDomain.ActivePeriod}
	if err := periodRepo.Create(period); err != nil {
		t.Fatalf("expected period seed, got %v", err)
	}

	users := []*usersDomain.User{
		{ID: 1, Name: "Admin", Email: "admin@example.com", Password: "secret123", GlobalRole: usersDomain.RoleAdmin},
		{ID: 10, Name: "Professor One", Email: "prof1@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
		{ID: 11, Name: "Professor Two", Email: "prof2@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
	}
	for _, user := range users {
		if err := userRepo.Create(user); err != nil {
			t.Fatalf("expected user seed, got %v", err)
		}
	}

	workspaces := []*workspacesDomain.Workspace{
		{ID: 1, PeriodID: 1, UserID: 10, Name: "Workspace 1", Type: workspacesDomain.CourseType, InitialDate: "2027-01-05", FinalDate: "2027-04-26", Observations: "obs", State: workspacesDomain.ActiveState},
		{ID: 2, PeriodID: 1, UserID: 11, Name: "Workspace 2", Type: workspacesDomain.CourseType, InitialDate: "2027-01-05", FinalDate: "2027-04-26", Observations: "obs", State: workspacesDomain.ActiveState},
	}
	for _, workspace := range workspaces {
		if err := workspaceRepo.Create(workspace); err != nil {
			t.Fatalf("expected workspace seed, got %v", err)
		}
	}
}

func withWorkspaceCurrentUser(user authDomain.AuthenticatedUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("current_user", user)
		c.Next()
	}
}

func TestProfessorCannotUpdateForeignWorkspace(t *testing.T) {
	handler := setupWorkspaceHandlerForVisibilityTests(t)

	router := gin.New()
	router.Use(withWorkspaceCurrentUser(authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor}))
	router.PUT("/workspaces/:id", handler.UpdateWorkspace)

	body := map[string]any{"period_id": 1, "name": "Updated", "type": "course", "initial_date": "2027-01-05", "final_date": "2027-04-26", "observations": "obs", "state": "active"}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/workspaces/2", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusForbidden, w.Code, w.Body.String())
	}
}

func TestAdminCanUpdateAnyWorkspace(t *testing.T) {
	handler := setupWorkspaceHandlerForVisibilityTests(t)

	router := gin.New()
	router.Use(withWorkspaceCurrentUser(authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin}))
	router.PUT("/workspaces/:id", handler.UpdateWorkspace)

	body := map[string]any{"period_id": 1, "name": "Updated By Admin", "type": "course", "initial_date": "2027-01-05", "final_date": "2027-04-26", "observations": "obs", "state": "active"}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/workspaces/2", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, w.Code, w.Body.String())
	}
}