package delivery

import (
	assignmentsInfrastructure "backend/internal/assignments/infrastructure"
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

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer) {
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

	createTask := application.NewCreateTask(repo, assignmentRepo, workspaceRepo, weeksRepo, nil)
	listTasks := application.NewListTasks(repo)
	getTaskByID := application.NewGetTaskByID(repo)
	updateTask := application.NewUpdateTask(repo, assignmentRepo, nil)
	deleteTask := application.NewDeleteTask(repo, nil)
	handler := NewTaskHandler(createTask, listTasks, getTaskByID, updateTask, deleteTask, assignmentRepo, workspaceRepo)

	tasks := r.Group("/tasks")
	{
		tasks.Use(authorizer.RequireAuthentication())

		taskOperators := tasks.Group("")
		taskOperators.Use(authorizer.RequireRoles(usersDomain.RoleMonitor, usersDomain.RoleAssistant, usersDomain.RoleProfessor, usersDomain.RoleAdmin))
		taskOperators.POST("", handler.CreateTask)
		taskOperators.GET("", handler.ListTasks)
		taskOperators.GET("/:id", handler.GetTaskByID)
		taskOperators.PUT("/:id", handler.UpdateTask)
		taskOperators.DELETE("/:id", handler.DeleteTask)
	}
}
