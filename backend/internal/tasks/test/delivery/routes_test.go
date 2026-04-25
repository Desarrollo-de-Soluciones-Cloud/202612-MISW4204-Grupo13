package delivery

import (
	"backend/internal/shared/database"
	tasksDelivery "backend/internal/tasks/delivery"
	usersDomain "backend/internal/users/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type noopAuthorizer struct{}

func (noopAuthorizer) RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func (noopAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesRegistersTasksEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	router := gin.New()
	api := router.Group("/api")
	tasksDelivery.SetupRoutes(api, noopAuthorizer{})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code == http.StatusNotFound {
		t.Fatal("expected tasks routes to be registered")
	}
}
