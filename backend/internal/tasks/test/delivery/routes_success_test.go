package delivery_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	sharedConfig "backend/internal/shared/config"
	sharedDB "backend/internal/shared/database/testsupport"
	tasksDelivery "backend/internal/tasks/delivery"
	tasksDomain "backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

const testTaskIDRoute = "/tasks/:id"

func TestSetupRoutesRegistersTaskRoutesWithBucket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("STORAGE_EMULATOR_HOST", "http://127.0.0.1:9090")
	sharedDB.SetupSQLiteDB(t, &tasksDomain.Task{}, &assignmentsDomain.Assignment{}, &workspacesDomain.Workspace{}, &weeksDomain.Week{})
	router := gin.New()

	tasksDelivery.SetupRoutes(router, taskRouteTestAuthorizer{}, &sharedConfig.Config{
		GCSBucketName:        "bucket",
		GCSAttachmentsPrefix: "attachments",
	})

	assertTaskRouteRegistered(t, router.Routes(), http.MethodPost, "/tasks")
	assertTaskRouteRegistered(t, router.Routes(), http.MethodGet, "/tasks")
	assertTaskRouteRegistered(t, router.Routes(), http.MethodGet, testTaskIDRoute)
	assertTaskRouteRegistered(t, router.Routes(), http.MethodGet, "/tasks/:id/attachments/:attachmentId/download")
	assertTaskRouteRegistered(t, router.Routes(), http.MethodPut, testTaskIDRoute)
	assertTaskRouteRegistered(t, router.Routes(), http.MethodDelete, testTaskIDRoute)
}

func assertTaskRouteRegistered(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("expected route %s %s to be registered", method, path)
}
