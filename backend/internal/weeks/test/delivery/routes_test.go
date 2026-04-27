package delivery_test

import (
	"backend/internal/shared/database"
	usersDomain "backend/internal/users/domain"
	"backend/internal/weeks/delivery"
	weeksDomain "backend/internal/weeks/domain"
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

	delivery.SetupRoutes(router, noopAuthorizer{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/periods/1", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Fatalf("expected route /weeks/periods/:periodId to be registered")
	}
}

func TestSetupRoutesRegistersWeekByNumberRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	if err := database.DB.AutoMigrate(&weeksDomain.Week{}); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	if err := database.DB.Create(&weeksDomain.Week{
		PeriodID:    1,
		Number:      1,
		InitialDate: "2026-01-12",
		FinalDate:   "2026-01-18",
	}).Error; err != nil {
		t.Fatalf("expected seed week, got %v", err)
	}

	router := gin.New()

	delivery.SetupRoutes(router, noopAuthorizer{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/1/periods/1", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Fatalf("expected route /weeks/:number/periods/:periodId to be registered")
	}
}
