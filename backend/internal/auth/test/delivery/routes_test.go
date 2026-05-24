package delivery_test

import (
	authDelivery "backend/internal/auth/delivery"
	"backend/internal/shared/config"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutesRegistersAuthEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := authDelivery.SetupRoutes(router, &config.Config{
		JWTSecret:            "secret",
		JWTExpirationMinutes: 60,
	})
	if handler == nil {
		t.Fatalf("expected auth handler")
	}

	routes := router.Routes()
	assertRouteRegistered(t, routes, http.MethodPost, "/auth/sign-in")
	assertRouteRegistered(t, routes, http.MethodGet, "/auth/me")
}

func assertRouteRegistered(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()

	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}

	t.Fatalf("expected route %s %s to be registered", method, path)
}
