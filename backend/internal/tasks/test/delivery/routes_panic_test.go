package delivery_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	sharedConfig "backend/internal/shared/config"
	sharedDB "backend/internal/shared/database/testsupport"
	tasksDelivery "backend/internal/tasks/delivery"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"testing"

	"github.com/gin-gonic/gin"
)

type taskRouteTestAuthorizer struct{}

func (taskRouteTestAuthorizer) RequireAuthentication() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func (taskRouteTestAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesPanicsWithoutBucketNameInTestFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sharedDB.SetupSQLiteDB(t, &tasksDomain.Task{}, &assignmentsDomain.Assignment{}, &workspacesDomain.Workspace{}, &weeksDomain.Week{})
	router := gin.New()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without bucket name")
		}
	}()

	tasksDelivery.SetupRoutes(router, taskRouteTestAuthorizer{}, &sharedConfig.Config{})
}
