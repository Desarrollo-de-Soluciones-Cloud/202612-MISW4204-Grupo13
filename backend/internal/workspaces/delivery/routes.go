package delivery

import (
	assignmentsInfra "backend/internal/assignments/infrastructure"
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

	// Initialize assignment repository for monitors and assistants
	assignmentRepo := assignmentsInfra.NewAssignmentRepository()

	createWorkspace := workspacesApp.NewCreateWorkspace(workspaceRepo, periodRepo, userRepo)
	listWorkspaces := workspacesApp.NewListWorkspaces(workspaceRepo)
	listWorkspacesByPeriod := workspacesApp.NewListWorkspacesByPeriod(workspaceRepo)
	getWorkspaceByID := workspacesApp.NewGetWorkspaceByID(workspaceRepo)
	updateWorkspace := workspacesApp.NewUpdateWorkspace(workspaceRepo, periodRepo, userRepo)
	deleteWorkspace := workspacesApp.NewDeleteWorkspace(workspaceRepo)
	closeWorkspace := workspacesApp.NewCloseWorkspace(workspaceRepo, userRepo)
	listWorkspaceMonitorsAndAssistants := workspacesApp.NewListWorkspaceMonitorsAndAssistants(workspaceRepo, assignmentRepo, userRepo)
	handler := NewWorkspaceHandler(WorkspaceHandlerUseCases{
		CreateWorkspace:                    createWorkspace,
		ListWorkspaces:                     listWorkspaces,
		ListWorkspacesByPeriod:             listWorkspacesByPeriod,
		GetWorkspaceByID:                   getWorkspaceByID,
		UpdateWorkspace:                    updateWorkspace,
		DeleteWorkspace:                    deleteWorkspace,
		CloseWorkspace:                     closeWorkspace,
		ListWorkspaceMonitorsAndAssistants: listWorkspaceMonitorsAndAssistants,
	})
	workspaces := r.Group("/workspaces")
	{
		workspaces.Use(authorizer.RequireAuthentication())

		adminAndProfessorWorkspaces := workspaces.Group("")
		adminAndProfessorWorkspaces.Use(authorizer.RequireRoles(usersDomain.RoleAdmin, usersDomain.RoleProfessor))
		adminAndProfessorWorkspaces.POST("", handler.CreateWorkspace)
		adminAndProfessorWorkspaces.GET("", handler.ListWorkspaces)
		adminAndProfessorWorkspaces.GET("/:id", handler.GetWorkspaceByID)
		adminAndProfessorWorkspaces.PUT("/:id", handler.UpdateWorkspace)
		adminAndProfessorWorkspaces.PATCH("/:id/close", handler.CloseWorkspace)

		professorWorkspaces := workspaces.Group("")
		professorWorkspaces.Use(authorizer.RequireRoles(usersDomain.RoleProfessor))
		professorWorkspaces.GET("/monitors-and-assistants/list", handler.ListWorkspaceMonitorsAndAssistants)
	}
}
