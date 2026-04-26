package delivery

import (
	"backend/internal/assignments/application"
	"backend/internal/assignments/infrastructure"
	periodsInfrastructure "backend/internal/periods/infrastructure"
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
	periodRepo := periodsInfrastructure.NewPeriodRepository()

	createAssignment := application.NewCreateAssignment(repo).WithRepositories(userRepo, workspaceRepo).WithPeriodRepository(periodRepo)
	getAssignmentByID := application.NewGetAssignmentByID(repo)
	listAllAssignments := application.NewListAllAssignments(repo)
	listAssignmentsByWorkspace := application.NewListAssignmentsByWorkspace(repo)
	listAssignmentsByUserID := application.NewListAssignmentsByUserID(repo)
	updateAssignment := application.NewUpdateAssignment(repo).WithWorkspaceRepository(workspaceRepo).WithPeriodRepository(periodRepo)

	handler := NewAssignmentHandler(createAssignment, getAssignmentByID, listAllAssignments, listAssignmentsByWorkspace, listAssignmentsByUserID, updateAssignment, workspaceRepo, userRepo)

	assignments := r.Group("/assignments")
	{
		assignments.Use(authorizer.RequireAuthentication())

		adminAndProfessorAssignments := assignments.Group("")
		adminAndProfessorAssignments.Use(authorizer.RequireRoles(usersDomain.RoleAdmin, usersDomain.RoleProfessor))
		adminAndProfessorAssignments.POST("", handler.CreateAssignment)
		adminAndProfessorAssignments.PUT("/:id", handler.UpdateAssignment)

		assignmentReaders := assignments.Group("")
		assignmentReaders.Use(authorizer.RequireRoles(usersDomain.RoleProfessor, usersDomain.RoleAdmin, usersDomain.RoleMonitor, usersDomain.RoleAssistant))
		assignmentReaders.GET("/:id", handler.GetAssignmentByID)
		assignmentReaders.GET("", handler.ListAssignments)
	}
}
