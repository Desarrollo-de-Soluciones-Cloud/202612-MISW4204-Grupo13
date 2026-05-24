package delivery_test

import (
	authDomain "backend/internal/auth/domain"
	deliverypkg "backend/internal/reports/delivery"
	usersDomain "backend/internal/users/domain"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newReportHandlerForTest() *deliverypkg.ReportHandler {
	return deliverypkg.NewReportHandler(nil, nil, nil, nil, nil, nil, nil, "")
}

func TestGenerateWeeklyReportsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/reports/weekly", nil)

	handler.GenerateWeeklyReports(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGenerateWeeklyReportsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly", bytes.NewBufferString(`{"workspace_id":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.GenerateWeeklyReports(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListReportsBadWorkspaceFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reports", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.ListReports(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDownloadReportBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "bad"}}
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.DownloadReport(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
