package delivery

import (
	"backend/internal/shared/database"
	usersDelivery "backend/internal/users/delivery"
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

func TestSetupRoutesRegistersUsersEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	if err := database.DB.AutoMigrate(&usersDomain.User{}); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	router := gin.New()
	api := router.Group("/api")
	usersDelivery.SetupRoutes(api, noopAuthorizer{})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/users", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code == http.StatusNotFound {
		t.Fatal("expected users routes to be registered")
	}
}
