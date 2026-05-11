package delivery

import (
	"context"
	assignmentsInfrastructure "backend/internal/assignments/infrastructure"
	sharedConfig "backend/internal/shared/config"
	sharedStorage "backend/internal/shared/storage"
	"backend/internal/tasks/application"
	"backend/internal/tasks/infrastructure"
	usersDomain "backend/internal/users/domain"
	weeksInfrastructure "backend/internal/weeks/infrastructure"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...usersDomain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer, cfg *sharedConfig.Config) {
	repo := infrastructure.NewTaskRepository()
	if err := repo.AutoMigrate(); err != nil {
		panic(err)
	}
	if err := repo.NormalizeLegacyStatuses(); err != nil {
		panic(err)
	}
	assignmentRepo := assignmentsInfrastructure.NewAssignmentRepository()
	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()
	weeksRepo := weeksInfrastructure.NewWeekRepository()
	taskFileStorage, err := sharedStorage.NewGCSStorage(context.Background(), cfg.GCSBucketName)
	if err != nil {
		panic(err)
	}

	createTask := application.NewCreateTask(repo, assignmentRepo, workspaceRepo, weeksRepo, nil)
	listTasks := application.NewListTasks(repo)
	getTaskByID := application.NewGetTaskByID(repo)
	updateTask := application.NewUpdateTask(repo, assignmentRepo, nil)
	setTaskAttachments := application.NewSetTaskAttachments(repo)
	deleteTask := application.NewDeleteTask(repo, nil)
	handler := NewTaskHandler(createTask, listTasks, getTaskByID, updateTask, setTaskAttachments, deleteTask, assignmentRepo, workspaceRepo, taskFileStorage, cfg.GCSAttachmentsPrefix)

	tasks := r.Group("/tasks")
	{
		tasks.Use(authorizer.RequireAuthentication())

		taskOperators := tasks.Group("")
		taskOperators.Use(authorizer.RequireRoles(usersDomain.RoleMonitor, usersDomain.RoleAssistant, usersDomain.RoleProfessor, usersDomain.RoleAdmin))
		taskOperators.POST("", handler.CreateTask)
		taskOperators.GET("", handler.ListTasks)
		taskOperators.GET("/:id", handler.GetTaskByID)
		taskOperators.GET("/:id/attachments/:attachmentId/download", handler.DownloadAttachment)
		taskOperators.PUT("/:id", handler.UpdateTask)
		taskOperators.DELETE("/:id", handler.DeleteTask)
	}
}
