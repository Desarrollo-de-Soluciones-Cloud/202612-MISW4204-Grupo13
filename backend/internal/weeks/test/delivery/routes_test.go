package delivery_test

import (
	"backend/internal/shared/database"
	"backend/internal/weeks/delivery"
	weeksDomain "backend/internal/weeks/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSetupRoutesRegistersWeekReadRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	if err := database.DB.AutoMigrate(&weeksDomain.Week{}); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	router := gin.New()

	delivery.SetupRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods/1/weeks", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Fatalf("expected route /periods/:periodId/weeks to be registered")
	}
}
