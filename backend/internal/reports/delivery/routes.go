package delivery

import (
	reportsApplication "backend/internal/reports/application"
	reportsInfrastructure "backend/internal/reports/infrastructure"
	usersDomain "backend/internal/users/domain"
	weeksInfrastructure "backend/internal/weeks/infrastructure"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...usersDomain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer) {
	reportRepo := reportsInfrastructure.NewReportRepository()
	if err := reportRepo.AutoMigrate(); err != nil {
		panic(err)
	}

	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()
	weekRepo := weeksInfrastructure.NewWeekRepository()
	assignmentReader := reportsInfrastructure.NewAssignmentReader()
	taskReader := reportsInfrastructure.NewTaskReader()
	pdfGenerator := reportsInfrastructure.NewPDFGenerator()

	generateWeeklyReports := reportsApplication.NewGenerateWeeklyReports(
		reportRepo,
		workspaceRepo,
		weekRepo,
		assignmentReader,
		taskReader,
		pdfGenerator,
		nil,
	)
	listReports := reportsApplication.NewListReports(reportRepo)
	getReportByID := reportsApplication.NewGetReportByID(reportRepo)

	handler := NewReportHandler(generateWeeklyReports, listReports, getReportByID, workspaceRepo)

	reports := r.Group("/reports")
	{
		reports.Use(authorizer.RequireAuthentication())

		reportOperators := reports.Group("")
		reportOperators.Use(authorizer.RequireRoles(usersDomain.RoleAdmin, usersDomain.RoleProfessor))
		reportOperators.POST("/weekly", handler.GenerateWeeklyReports)
		reportOperators.GET("", handler.ListReports)
		reportOperators.GET("/:id/download", handler.DownloadReport)
	}
}
