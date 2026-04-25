package delivery_test

import (
	authDomain "backend/internal/auth/domain"
	assignmentsDomain "backend/internal/assignments/domain"
	assignmentsInfrastructure "backend/internal/assignments/infrastructure"
	periodsDomain "backend/internal/periods/domain"
	periodsInfrastructure "backend/internal/periods/infrastructure"
	reportsApplication "backend/internal/reports/application"
	reportsDelivery "backend/internal/reports/delivery"
	reportsInfrastructure "backend/internal/reports/infrastructure"
	"backend/internal/shared/database"
	tasksDomain "backend/internal/tasks/domain"
	tasksInfrastructure "backend/internal/tasks/infrastructure"
	usersDomain "backend/internal/users/domain"
	usersInfrastructure "backend/internal/users/infrastructure"
	weeksDomain "backend/internal/weeks/domain"
	weeksInfrastructure "backend/internal/weeks/infrastructure"
	workspacesDomain "backend/internal/workspaces/domain"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	workspaceInitialDate      = "2027-01-04"
	workspaceFinalDate        = "2027-04-25"
	reportsWeeklyPath         = "/reports/weekly"
	contentTypeHeader         = "Content-Type"
	applicationJSON           = "application/json"
	expectedStatusErrorFormat = "expected status %d, got %d (body: %s)"
)

func setupReportHandlerForDeliveryTests(t *testing.T) *reportsDelivery.ReportHandler {
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
	weekRepo := weeksInfrastructure.NewWeekRepository()
	assignmentRepo := assignmentsInfrastructure.NewAssignmentRepository()
	taskRepo := tasksInfrastructure.NewTaskRepository()
	reportRepo := reportsInfrastructure.NewReportRepository()

	if err := periodRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected periods automigrate, got %v", err)
	}
	if err := userRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected users automigrate, got %v", err)
	}
	if err := workspaceRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected workspaces automigrate, got %v", err)
	}
	if err := weekRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected weeks automigrate, got %v", err)
	}
	if err := assignmentRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected assignments automigrate, got %v", err)
	}
	if err := taskRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected tasks automigrate, got %v", err)
	}
	if err := reportRepo.AutoMigrate(); err != nil {
		t.Fatalf("expected reports automigrate, got %v", err)
	}

	seedReportData(t, periodRepo, userRepo, workspaceRepo, weekRepo, assignmentRepo, taskRepo)

	reportsDir := filepath.Join(t.TempDir(), "reports")
	generateWeeklyReports := reportsApplication.NewGenerateWeeklyReports(
		reportRepo,
		workspaceRepo,
		weekRepo,
		reportsInfrastructure.NewAssignmentReader(),
		reportsInfrastructure.NewTaskReader(),
		reportsInfrastructure.NewPDFGenerator(),
		&reportsApplication.GenerateWeeklyReportsOptions{
			ReportsStorageDir: reportsDir,
			Now:               func() time.Time { return time.Date(2027, time.January, 1, 10, 0, 0, 0, time.UTC) },
		},
	)
	listReports := reportsApplication.NewListReports(reportRepo)
	getReportByID := reportsApplication.NewGetReportByID(reportRepo)

	return reportsDelivery.NewReportHandler(generateWeeklyReports, listReports, getReportByID, workspaceRepo)
}

func seedReportData(
	t *testing.T,
	periodRepo *periodsInfrastructure.PeriodRepository,
	userRepo *usersInfrastructure.UserRepository,
	workspaceRepo *workspacesInfrastructure.WorkspaceRepository,
	weekRepo *weeksInfrastructure.WeekRepository,
	assignmentRepo *assignmentsInfrastructure.AssignmentRepository,
	taskRepo *tasksInfrastructure.TaskRepository,
) {
	t.Helper()

	period := &periodsDomain.Period{
		ID:                   1,
		Name:                 "2027-10",
		InitialDate:          workspaceInitialDate,
		FinalDate:            workspaceFinalDate,
		InscriptionFinalDate: "2026-12-31",
		WeeksCount:           16,
		PeriodState:          periodsDomain.ActivePeriod,
	}
	if err := periodRepo.Create(period); err != nil {
		t.Fatalf("expected period seed, got %v", err)
	}

	users := []*usersDomain.User{
		{ID: 1, Name: "Admin", Email: "admin@example.com", Password: "secret123", GlobalRole: usersDomain.RoleAdmin},
		{ID: 10, Name: "Professor 1", Email: "prof1@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
		{ID: 11, Name: "Professor 2", Email: "prof2@example.com", Password: "secret123", GlobalRole: usersDomain.RoleProfessor},
		{ID: 100, Name: "Monitor", Email: "monitor@example.com", Password: "secret123", GlobalRole: usersDomain.RoleMonitor},
	}
	for _, user := range users {
		if err := userRepo.Create(user); err != nil {
			t.Fatalf("expected user seed, got %v", err)
		}
	}

	workspaces := []*workspacesDomain.Workspace{
		{ID: 1, PeriodID: 1, UserID: 10, Name: "Workspace A", Type: workspacesDomain.CourseType, InitialDate: workspaceInitialDate, FinalDate: workspaceFinalDate, Observations: "obs", State: workspacesDomain.ActiveState},
		{ID: 2, PeriodID: 1, UserID: 11, Name: "Workspace B", Type: workspacesDomain.CourseType, InitialDate: workspaceInitialDate, FinalDate: workspaceFinalDate, Observations: "obs", State: workspacesDomain.ActiveState},
	}
	for _, workspace := range workspaces {
		if err := workspaceRepo.Create(workspace); err != nil {
			t.Fatalf("expected workspace seed, got %v", err)
		}
	}

	week := &weeksDomain.Week{ID: 1, PeriodID: 1, Number: 1, InitialDate: "2027-01-04", FinalDate: "2027-01-10"}
	if err := weekRepo.CreateMany([]weeksDomain.Week{*week}); err != nil {
		t.Fatalf("expected week seed, got %v", err)
	}

	assignment := &assignmentsDomain.Assignment{ID: 1, UserID: 100, WorkspaceID: 1, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 8}
	if err := assignmentRepo.Create(assignment); err != nil {
		t.Fatalf("expected assignment seed, got %v", err)
	}

	weekID := uint(1)
	task := &tasksDomain.Task{
		ID:            1,
		UserID:        100,
		AssignmentID:  1,
		WeekID:        &weekID,
		Title:         "Implementar endpoint",
		Description:   "Se implementa endpoint de tareas",
		Status:        tasksDomain.TaskStatusFinalizado,
		SpentHours:    5,
		Observations:  "ok",
		WeekStartDate: time.Date(2027, time.January, 4, 0, 0, 0, 0, time.UTC),
		Late:          false,
	}
	if err := taskRepo.Create(task); err != nil {
		t.Fatalf("expected task seed, got %v", err)
	}
}

func withReportCurrentUser(user authDomain.AuthenticatedUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("current_user", user)
		c.Next()
	}
}

func TestProfessorCannotGenerateWeeklyReportForForeignWorkspace(t *testing.T) {
	handler := setupReportHandlerForDeliveryTests(t)

	router := gin.New()
	router.Use(withReportCurrentUser(authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor}))
	router.POST(reportsWeeklyPath, handler.GenerateWeeklyReports)

	body := map[string]any{"workspace_id": 2, "week_id": 1}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, reportsWeeklyPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(contentTypeHeader, applicationJSON)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf(expectedStatusErrorFormat, http.StatusForbidden, w.Code, w.Body.String())
	}
}

func TestAdminCanGenerateWeeklyReport(t *testing.T) {
	handler := setupReportHandlerForDeliveryTests(t)

	router := gin.New()
	router.Use(withReportCurrentUser(authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin}))
	router.POST(reportsWeeklyPath, handler.GenerateWeeklyReports)

	body := map[string]any{"workspace_id": 1, "week_id": 1}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, reportsWeeklyPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(contentTypeHeader, applicationJSON)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf(expectedStatusErrorFormat, http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestProfessorCanGenerateWeeklyReportForOwnWorkspace(t *testing.T) {
	handler := setupReportHandlerForDeliveryTests(t)

	router := gin.New()
	router.Use(withReportCurrentUser(authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor}))
	router.POST(reportsWeeklyPath, handler.GenerateWeeklyReports)

	body := map[string]any{"workspace_id": 1, "week_id": 1}
	bodyJSON, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, reportsWeeklyPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(contentTypeHeader, applicationJSON)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf(expectedStatusErrorFormat, http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestDownloadNonExistentReportReturnsNotFound(t *testing.T) {
	handler := setupReportHandlerForDeliveryTests(t)

	router := gin.New()
	router.Use(withReportCurrentUser(authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin}))
	router.GET("/reports/:id/download", handler.DownloadReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/reports/999/download", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf(expectedStatusErrorFormat, http.StatusNotFound, w.Code, w.Body.String())
	}
}
