package delivery_test

import (
	usersDomain "backend/internal/users/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateWorkspaceUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/workspaces", nil)

	handler.CreateWorkspace(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateWorkspaceProfessorForbiddenForOtherOwner(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]any{
		"period_id": 1, "user_id": 99, "name": "Algorithms", "type": "course",
		"initial_date": "2026-06-02", "final_date": "2026-06-30", "observations": "obs", "state": "active",
	})
	req := httptest.NewRequest(http.MethodPost, "/workspaces", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authenticatedUser(10, usersDomain.RoleProfessor))

	handler.CreateWorkspace(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestListWorkspacesProfessorFiltersOwnRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/workspaces", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authenticatedUser(10, usersDomain.RoleProfessor))

	handler.ListWorkspaces(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var response struct {
		Workspaces []any `json:"workspaces"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected json, got %v", err)
	}
	if len(response.Workspaces) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(response.Workspaces))
	}
}

func TestGetWorkspaceByIDBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/workspaces/bad", nil)
	c.Params = gin.Params{{Key: "id", Value: "bad"}}
	c.Set("current_user", authenticatedUser(1, usersDomain.RoleAdmin))

	handler.GetWorkspaceByID(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateWorkspaceBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/workspaces/1", bytes.NewBufferString(`{"name":1}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("current_user", authenticatedUser(1, usersDomain.RoleAdmin))

	handler.UpdateWorkspace(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCloseWorkspaceSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPatch, "/workspaces/1/close", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("current_user", authenticatedUser(1, usersDomain.RoleAdmin))

	handler.CloseWorkspace(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListWorkspaceMonitorsAndAssistantsProfessorOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/workspaces/monitors-and-assistants/list", nil)
	c.Set("current_user", authenticatedUser(20, usersDomain.RoleMonitor))

	handler.ListWorkspaceMonitorsAndAssistants(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestListWorkspaceMonitorsAndAssistantsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWorkspaceHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/workspaces/monitors-and-assistants/list", nil)
	c.Set("current_user", authenticatedUser(10, usersDomain.RoleProfessor))

	handler.ListWorkspaceMonitorsAndAssistants(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
