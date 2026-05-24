package delivery_test

import (
	sharedDB "backend/internal/shared/database/testsupport"
	usersDomain "backend/internal/users/domain"
	weeksDelivery "backend/internal/weeks/delivery"
	weeksDomain "backend/internal/weeks/domain"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type weeksRouteTestAuthorizer struct{}

func (weeksRouteTestAuthorizer) RequireAuthentication() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func (weeksRouteTestAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesRegistersWeekRoutesSQLite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sharedDB.SetupSQLiteDB(t, &weeksDomain.Week{})
	router := gin.New()
	weeksDelivery.SetupRoutes(router, weeksRouteTestAuthorizer{})

	assertWeeksRoute(t, router.Routes(), http.MethodGet, "/weeks/periods/:periodId")
	assertWeeksRoute(t, router.Routes(), http.MethodGet, "/weeks/:number/periods/:periodId")
}

func assertWeeksRoute(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("expected route %s %s", method, path)
}
