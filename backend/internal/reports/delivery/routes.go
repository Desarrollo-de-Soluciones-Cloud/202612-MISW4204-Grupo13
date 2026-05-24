package delivery

import (
	"context"

	reportsApplication "backend/internal/reports/application"
	reportsInfrastructure "backend/internal/reports/infrastructure"
	sharedConfig "backend/internal/shared/config"
	sharedStorage "backend/internal/shared/storage"
	usersDomain "backend/internal/users/domain"
	usersInfrastructure "backend/internal/users/infrastructure"
	weeksInfrastructure "backend/internal/weeks/infrastructure"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...usersDomain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer, cfg *sharedConfig.Config) {
	reportRepo := reportsInfrastructure.NewReportRepository()
	if err := reportRepo.AutoMigrate(); err != nil {
		panic(err)
	}

	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()
	weekRepo := weeksInfrastructure.NewWeekRepository()
	userRepo := usersInfrastructure.NewUserRepository()

	assignmentReader := reportsInfrastructure.NewAssignmentReader()
	taskReader := reportsInfrastructure.NewTaskReader()
	pdfGenerator := reportsInfrastructure.NewPDFGenerator()

	aiReportGenerator, err := reportsInfrastructure.NewVertexAIReportGenerator(
		cfg.GCPProjectID,
		cfg.GCPLocation,
		cfg.VertexAIModel,
	)
	if err != nil {
		panic(err)
	}

	reportFileStorage, err := sharedStorage.NewGCSStorage(
		context.Background(),
		cfg.GCSBucketName,
	)
	if err != nil {
		panic(err)
	}

	generateWeeklyReports := reportsApplication.NewGenerateWeeklyReports(
		reportRepo,
		workspaceRepo,
		weekRepo,
		assignmentReader,
		taskReader,
		userRepo,
		pdfGenerator,
		aiReportGenerator,
		reportFileStorage,
		&reportsApplication.GenerateWeeklyReportsOptions{
			ReportsGCSPrefix: cfg.GCSReportsPrefix,
		},
	)

	processWeeklyReportJob := reportsApplication.NewProcessWeeklyReportJob(
		reportRepo,
		workspaceRepo,
		weekRepo,
		assignmentReader,
		taskReader,
		userRepo,
		pdfGenerator,
		aiReportGenerator,
		reportFileStorage,
		&reportsApplication.GenerateWeeklyReportsOptions{
			ReportsGCSPrefix: cfg.GCSReportsPrefix,
		},
	)

	var queueWeeklyReports *reportsApplication.QueueWeeklyReports
	if cfg.GCPProjectID != "" && cfg.ReportsPubSubTopic != "" {
		publisher, err := reportsInfrastructure.NewPubSubReportJobPublisher(
			context.Background(),
			cfg.GCPProjectID,
			cfg.ReportsPubSubTopic,
		)
		if err != nil {
			panic(err)
		}

		queueWeeklyReports = reportsApplication.NewQueueWeeklyReports(
			workspaceRepo,
			weekRepo,
			assignmentReader,
			publisher,
		)
	}

	listReports := reportsApplication.NewListReports(
		reportRepo,
		workspaceRepo,
		weekRepo,
		userRepo,
	)

	getReportByID := reportsApplication.NewGetReportByID(reportRepo)

	handler := NewReportHandler(
		generateWeeklyReports,
		queueWeeklyReports,
		processWeeklyReportJob,
		listReports,
		getReportByID,
		workspaceRepo,
		reportFileStorage,
		cfg.PubSubPushAuthToken,
	)

	reports := r.Group("/reports")
	reports.POST("/weekly/process", handler.ProcessWeeklyReportJob)

	reports.Use(authorizer.RequireAuthentication())

	reportOperators := reports.Group("")
	reportOperators.Use(authorizer.RequireRoles(usersDomain.RoleAdmin, usersDomain.RoleProfessor))

	reportOperators.POST("/weekly", handler.GenerateWeeklyReports)
	reportOperators.GET("", handler.ListReports)
	reportOperators.GET("/:id/download", handler.DownloadReport)
}
