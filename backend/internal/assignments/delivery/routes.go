package delivery

import (
	"backend/internal/assignments/application"
	"backend/internal/assignments/infrastructure"
	usersDomain "backend/internal/users/domain"
	usersInfrastructure "backend/internal/users/infrastructure"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...usersDomain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer) {
	repo := infrastructure.NewAssignmentRepository()
	repo.AutoMigrate()
	userRepo := usersInfrastructure.NewUserRepository()
	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()

	createAssignment := application.NewCreateAssignment(repo).WithRepositories(userRepo, workspaceRepo)
	getAssignmentByID := application.NewGetAssignmentByID(repo)
	listAssignmentsByUserID := application.NewListAssignmentsByUserID(repo)
	updateAssignment := application.NewUpdateAssignment(repo).WithWorkspaceRepository(workspaceRepo)

	handler := NewAssignmentHandler(createAssignment, getAssignmentByID, listAssignmentsByUserID, updateAssignment, workspaceRepo, userRepo)

	assignments := r.Group("/assignments")
	{
		assignments.Use(authorizer.RequireAuthentication())

		adminAndProfessorAssignments := assignments.Group("")
		adminAndProfessorAssignments.Use(authorizer.RequireRoles(usersDomain.RoleAdmin, usersDomain.RoleProfessor))
		adminAndProfessorAssignments.POST("", handler.CreateAssignment)

		assignmentReaders := assignments.Group("")
		assignmentReaders.Use(authorizer.RequireRoles(usersDomain.RoleProfessor, usersDomain.RoleAdmin, usersDomain.RoleMonitor, usersDomain.RoleAssistant))
		assignmentReaders.GET("/:id", handler.GetAssignmentByID)
		assignmentReaders.GET("", handler.ListAssignmentsByUserID)

		adminAssignments := assignments.Group("")
		adminAssignments.Use(authorizer.RequireRoles(usersDomain.RoleAdmin))
		adminAssignments.PUT("/:id/admin", handler.UpdateAssignment)
	}
}
