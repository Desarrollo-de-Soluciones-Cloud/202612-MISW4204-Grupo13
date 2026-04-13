package delivery

import (
	periodInfra "backend/internal/periods/infrastructure"
	usersInfra "backend/internal/users/infrastructure"
	workspacesApp "backend/internal/workspaces/application"
	"backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
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
		workspaces.POST("", handler.CreateWorkspace)
		workspaces.GET("", handler.ListWorkspaces)
		workspaces.GET("/:id", handler.GetWorkspaceByID)
		workspaces.PUT("/:id", handler.UpdateWorkspace)
		workspaces.DELETE("/:id", handler.DeleteWorkspace)
	}
}
