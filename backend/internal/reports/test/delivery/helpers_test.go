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

func TestGenerateWeeklyReportsMissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewReportHandler(nil, nil, nil, nil, nil, nil, nil, "")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly", bytes.NewBufferString(`{}`))
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
