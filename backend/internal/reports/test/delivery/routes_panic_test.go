package delivery_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDelivery "backend/internal/reports/delivery"
	reportsDomain "backend/internal/reports/domain"
	sharedConfig "backend/internal/shared/config"
	sharedDB "backend/internal/shared/database/testsupport"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"testing"

	"github.com/gin-gonic/gin"
)

type reportRouteTestAuthorizer struct{}

func (reportRouteTestAuthorizer) RequireAuthentication() gin.HandlerFunc { return func(c *gin.Context) { c.Next() } }
func (reportRouteTestAuthorizer) RequireRoles(...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}

func TestSetupRoutesPanicsWithoutGCPProjectInTestFolder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sharedDB.SetupSQLiteDB(t, &reportsDomain.Report{}, &workspacesDomain.Workspace{}, &weeksDomain.Week{}, &usersDomain.User{}, &assignmentsDomain.Assignment{})
	router := gin.New()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without gcp project")
		}
	}()

	reportsDelivery.SetupRoutes(router, reportRouteTestAuthorizer{}, &sharedConfig.Config{
		GCSBucketName:    "bucket",
		GCSReportsPrefix: "reports",
	})
}
