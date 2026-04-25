package delivery

import (
	periodInfra "backend/internal/periods/infrastructure"
	usersDomain "backend/internal/users/domain"
	usersInfra "backend/internal/users/infrastructure"
	workspacesApp "backend/internal/workspaces/application"
	"backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...usersDomain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer) {
	// Initialize workspace repository
	workspaceRepo := infrastructure.NewWorkspaceRepository()
	workspaceRepo.AutoMigrate()

	// Initialize period repository for workspace period validation
	periodRepo := periodInfra.NewPeriodRepository()

	// Initialize user repository for workspace user validation
	userRepo := usersInfra.NewUserRepository()

	createWorkspace := workspacesApp.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	listWorkspaces := workspacesApp.NewListWorkspaces(workspaceRepo)
	listWorkspacesByPeriod := workspacesApp.NewListWorkspacesByPeriod(workspaceRepo)
	getWorkspaceByID := workspacesApp.NewGetWorkspaceByID(workspaceRepo)
	updateWorkspace := workspacesApp.NewUpdateWorkspace(workspaceRepo, periodRepo, userRepo)
	deleteWorkspace := workspacesApp.NewDeleteWorkspace(workspaceRepo)
	handler := NewWorkspaceHandler(createWorkspace, listWorkspaces, listWorkspacesByPeriod, getWorkspaceByID, updateWorkspace, deleteWorkspace)
	workspaces := r.Group("/workspaces")
	{
		workspaces.Use(authorizer.RequireAuthentication())

		adminAndProfessorWorkspaces := workspaces.Group("")
		adminAndProfessorWorkspaces.Use(authorizer.RequireRoles(usersDomain.RoleAdmin, usersDomain.RoleProfessor))
		adminAndProfessorWorkspaces.POST("", handler.CreateWorkspace)
		adminAndProfessorWorkspaces.GET("", handler.ListWorkspaces)
		adminAndProfessorWorkspaces.GET("/:id", handler.GetWorkspaceByID)
		adminAndProfessorWorkspaces.PUT("/:id", handler.UpdateWorkspace)
		adminAndProfessorWorkspaces.DELETE("/:id", handler.DeleteWorkspace)
	}
}
