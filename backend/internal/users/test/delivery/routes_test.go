package delivery_test

import (
	sharedDB "backend/internal/shared/database/testsupport"
	usersDelivery "backend/internal/users/delivery"
	usersDomain "backend/internal/users/domain"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type userRouteTestAuthorizer struct{}

func (userRouteTestAuthorizer) RequireAuthentication() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func (userRouteTestAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesRegistersUserRoutesSQLite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sharedDB.SetupSQLiteDB(t, &usersDomain.User{})
	router := gin.New()
	usersDelivery.SetupRoutes(router, userRouteTestAuthorizer{})

	assertUserRoute(t, router.Routes(), http.MethodPost, "/users")
	assertUserRoute(t, router.Routes(), http.MethodPut, "/users/:id")
	assertUserRoute(t, router.Routes(), http.MethodPatch, "/users/:id/role")
	assertUserRoute(t, router.Routes(), http.MethodGet, "/users")
	assertUserRoute(t, router.Routes(), http.MethodGet, "/users/:id")
}

func assertUserRoute(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("expected route %s %s", method, path)
}
