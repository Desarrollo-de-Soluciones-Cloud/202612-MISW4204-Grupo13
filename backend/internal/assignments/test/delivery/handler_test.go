package delivery_test

import (
	assignmentsApplication "backend/internal/assignments/application"
	assignmentsDelivery "backend/internal/assignments/delivery"
	assignmentsDomain "backend/internal/assignments/domain"
	authDomain "backend/internal/auth/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type assignmentRepoStub struct {
	assignments []assignmentsDomain.Assignment
}

func (r *assignmentRepoStub) Create(assignment *assignmentsDomain.Assignment) error {
	assignment.ID = uint(len(r.assignments) + 1)
	r.assignments = append(r.assignments, *assignment)
	return nil
}

func (r *assignmentRepoStub) FindByID(id uint) (*assignmentsDomain.Assignment, error) {
	for _, assignment := range r.assignments {
		if assignment.ID == id {
			copy := assignment
			return &copy, nil
		}
	}
	return nil, assignmentsDomain.ErrAssignmentNotFound
}

func (r *assignmentRepoStub) FindAll() ([]assignmentsDomain.Assignment, error) {
	result := make([]assignmentsDomain.Assignment, len(r.assignments))
	copy(result, r.assignments)
	return result, nil
}

func (r *assignmentRepoStub) FindAllByUserID(userID uint) ([]assignmentsDomain.Assignment, error) {
	result := make([]assignmentsDomain.Assignment, 0)
	for _, assignment := range r.assignments {
		if assignment.UserID == userID {
			result = append(result, assignment)
		}
	}
	return result, nil
}

func (r *assignmentRepoStub) FindByWorkspaceUserID(workspaceUserID uint) ([]assignmentsDomain.Assignment, error) {
	return r.FindAll()
}

func (r *assignmentRepoStub) FindByWorkspaceIDsAndRoles(workspaceIDs []uint, roles []assignmentsDomain.AssignmentRole) ([]assignmentsDomain.Assignment, error) {
	return r.FindAll()
}

func (r *assignmentRepoStub) SumWeeklyHoursByUserAndRole(userID uint, role assignmentsDomain.AssignmentRole) (int, error) {
	total := 0
	for _, assignment := range r.assignments {
		if assignment.UserID == userID && assignment.Role == role {
			total += assignment.WeeklyHours
		}
	}
	return total, nil
}

func (r *assignmentRepoStub) CountAssignmentsByUserAndRole(userID uint, role assignmentsDomain.AssignmentRole) (int, error) {
	total := 0
	for _, assignment := range r.assignments {
		if assignment.UserID == userID && assignment.Role == role {
			total++
		}
	}
	return total, nil
}

func (r *assignmentRepoStub) Update(assignment *assignmentsDomain.Assignment) error {
	for i, item := range r.assignments {
		if item.ID == assignment.ID {
			r.assignments[i] = *assignment
			return nil
		}
	}
	return assignmentsDomain.ErrAssignmentNotFound
}

type assignmentWorkspaceReaderStub struct {
	workspaces map[uint]*workspacesDomain.Workspace
}

func (r *assignmentWorkspaceReaderStub) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	workspace, ok := r.workspaces[id]
	if !ok {
		return nil, workspacesDomain.ErrWorkspaceNotFound
	}
	return workspace, nil
}

type assignmentUserReaderStub struct {
	users map[uint]*usersDomain.User
}

func (r *assignmentUserReaderStub) FindByID(id uint) (*usersDomain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, usersDomain.ErrUserNotFound
	}
	return user, nil
}

type assignmentPeriodRepoStub struct{}

func (r *assignmentPeriodRepoStub) FindByID(id uint) (*workspacesDomain.Workspace, error) { return nil, nil }

func newAssignmentHandlerForTest() *assignmentsDelivery.AssignmentHandler {
	repo := &assignmentRepoStub{
		assignments: []assignmentsDomain.Assignment{
			{ID: 1, UserID: 2, WorkspaceID: 7, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 8},
		},
	}
	workspaceReader := &assignmentWorkspaceReaderStub{
		workspaces: map[uint]*workspacesDomain.Workspace{
			7: {ID: 7, UserID: 10, PeriodID: 1, Name: "Proyecto", Type: workspacesDomain.ProjectType, InitialDate: "2026-01-20", FinalDate: "2026-05-20", Observations: "obs", State: workspacesDomain.ActiveState},
		},
	}
	userReader := &assignmentUserReaderStub{
		users: map[uint]*usersDomain.User{
			2: {ID: 2, Name: "Mon", Email: "mon@example.com", GlobalRole: usersDomain.RoleMonitor},
			3: {ID: 3, Name: "Asis", Email: "asis@example.com", GlobalRole: usersDomain.RoleAssistant},
		},
	}

	updateUseCase := assignmentsApplication.NewUpdateAssignment(repo)
	updateUseCase = updateUseCase.WithWorkspaceRepository(workspaceReader)

	return assignmentsDelivery.NewAssignmentHandler(
		assignmentsApplication.NewCreateAssignment(repo).WithRepositories(userReader, workspaceReader),
		assignmentsApplication.NewGetAssignmentByID(repo),
		assignmentsApplication.NewListAllAssignments(repo),
		assignmentsApplication.NewListAssignmentsByWorkspace(repo),
		assignmentsApplication.NewListAssignmentsByUserID(repo),
		updateUseCase,
		assignmentsDelivery.AssignmentHandlerDependencies{
			WorkspaceReader: workspaceReader,
			UserReader:      userReader,
		},
	)
}

func authenticatedUser(id uint, role usersDomain.UserRole) authDomain.AuthenticatedUser {
	return authDomain.AuthenticatedUser{
		ID:         id,
		Name:       "User",
		Email:      "user@example.com",
		GlobalRole: role,
	}
}

func TestCreateAssignmentUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAssignmentHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/assignments", nil)

	handler.CreateAssignment(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateAssignmentSuccessForProfessor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAssignmentHandlerForTest()
	body, _ := json.Marshal(map[string]any{
		"user_id":       3,
		"workspace_id":  7,
		"role":          "assistant",
		"weekly_hours":  6,
	})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/assignments", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("current_user", authenticatedUser(10, usersDomain.RoleProfessor))

	handler.CreateAssignment(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestGetAssignmentByIDForbiddenForOperationalUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAssignmentHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/assignments/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("current_user", authenticatedUser(99, usersDomain.RoleMonitor))

	handler.GetAssignmentByID(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestListAssignmentsForAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAssignmentHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/assignments", nil)
	c.Set("current_user", authenticatedUser(1, usersDomain.RoleAdmin))

	handler.ListAssignments(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateAssignmentRejectsProfessorWeeklyHourChange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAssignmentHandlerForTest()
	body, _ := json.Marshal(map[string]any{
		"role":          "monitor",
		"weekly_hours":  10,
	})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/assignments/1", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("current_user", authenticatedUser(10, usersDomain.RoleProfessor))

	handler.UpdateAssignment(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
