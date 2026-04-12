package delivery

import (
	authDelivery "backend/internal/auth/delivery"
	"backend/internal/shared/config"
	"backend/internal/shared/database"
	usersDomain "backend/internal/users/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSetupRoutesRegistersAuthEndpoints(t *testing.T) {
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

	handler := authDelivery.SetupRoutes(api, &config.Config{
		JWTSecret:            "secret",
		JWTExpirationMinutes: 60,
	})
	if handler == nil {
		t.Fatal("expected auth handler")
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/sign-in", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code == http.StatusNotFound {
		t.Fatal("expected auth route to exist")
	}
}
