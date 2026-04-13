package delivery

import (
	"backend/internal/shared/database"
	tasksDelivery "backend/internal/tasks/delivery"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSetupRoutesRegistersTasksEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	router := gin.New()
	api := router.Group("/api")
	tasksDelivery.SetupRoutes(api)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code == http.StatusNotFound {
		t.Fatal("expected tasks routes to be registered")
	}
}
