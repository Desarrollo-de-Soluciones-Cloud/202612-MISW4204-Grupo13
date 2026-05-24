package delivery_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	periodsDomain "backend/internal/periods/domain"
	sharedDB "backend/internal/shared/database/testsupport"
	usersDomain "backend/internal/users/domain"
	workspacesDelivery "backend/internal/workspaces/delivery"
	workspacesDomain "backend/internal/workspaces/domain"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type workspaceRouteTestAuthorizer struct{}

func (workspaceRouteTestAuthorizer) RequireAuthentication() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func (workspaceRouteTestAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesRegistersWorkspaceRoutesSQLite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sharedDB.SetupSQLiteDB(t, &workspacesDomain.Workspace{}, &periodsDomain.Period{}, &usersDomain.User{}, &assignmentsDomain.Assignment{})
	router := gin.New()
	workspacesDelivery.SetupRoutes(router, workspaceRouteTestAuthorizer{})

	assertWorkspaceRoute(t, router.Routes(), http.MethodPost, "/workspaces")
	assertWorkspaceRoute(t, router.Routes(), http.MethodGet, "/workspaces")
	assertWorkspaceRoute(t, router.Routes(), http.MethodGet, "/workspaces/:id")
	assertWorkspaceRoute(t, router.Routes(), http.MethodPut, "/workspaces/:id")
	assertWorkspaceRoute(t, router.Routes(), http.MethodPatch, "/workspaces/:id/close")
	assertWorkspaceRoute(t, router.Routes(), http.MethodGet, "/workspaces/monitors-and-assistants/list")
}

func assertWorkspaceRoute(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("expected route %s %s", method, path)
}
