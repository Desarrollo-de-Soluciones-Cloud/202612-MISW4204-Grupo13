package delivery_test

import (
	periodsDelivery "backend/internal/periods/delivery"
	periodsDomain "backend/internal/periods/domain"
	sharedDB "backend/internal/shared/database/testsupport"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type periodRouteTestAuthorizer struct{}

func (periodRouteTestAuthorizer) RequireAuthentication() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func (periodRouteTestAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesRegistersPeriodRoutesSQLite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sharedDB.SetupSQLiteDB(t, &periodsDomain.Period{}, &weeksDomain.Week{})
	router := gin.New()
	periodsDelivery.SetupRoutes(router, periodRouteTestAuthorizer{})

	assertPeriodRoute(t, router.Routes(), http.MethodPost, "/periods")
	assertPeriodRoute(t, router.Routes(), http.MethodPatch, "/periods/:id")
	assertPeriodRoute(t, router.Routes(), http.MethodPatch, "/periods/:id/close")
	assertPeriodRoute(t, router.Routes(), http.MethodGet, "/periods")
}

func assertPeriodRoute(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("expected route %s %s", method, path)
}
