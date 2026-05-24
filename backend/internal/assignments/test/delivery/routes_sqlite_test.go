package delivery_test

import (
	assignmentsDelivery "backend/internal/assignments/delivery"
	assignmentsDomain "backend/internal/assignments/domain"
	periodsDomain "backend/internal/periods/domain"
	sharedDB "backend/internal/shared/database/testsupport"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type assignmentRouteTestAuthorizer struct{}

func (assignmentRouteTestAuthorizer) RequireAuthentication() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func (assignmentRouteTestAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesRegistersAssignmentRoutesSQLite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sharedDB.SetupSQLiteDB(t, &assignmentsDomain.Assignment{}, &usersDomain.User{}, &workspacesDomain.Workspace{}, &periodsDomain.Period{})
	router := gin.New()
	assignmentsDelivery.SetupRoutes(router, assignmentRouteTestAuthorizer{})

	assertAssignmentRoute(t, router.Routes(), http.MethodPost, "/assignments")
	assertAssignmentRoute(t, router.Routes(), http.MethodPut, "/assignments/:id")
	assertAssignmentRoute(t, router.Routes(), http.MethodGet, "/assignments/:id")
	assertAssignmentRoute(t, router.Routes(), http.MethodGet, "/assignments")
}

func assertAssignmentRoute(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("expected route %s %s", method, path)
}
