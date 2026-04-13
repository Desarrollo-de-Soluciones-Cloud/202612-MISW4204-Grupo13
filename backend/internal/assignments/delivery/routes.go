package delivery

import (
	"backend/internal/assignments/application"
	"backend/internal/assignments/infrastructure"
	usersInfrastructure "backend/internal/users/infrastructure"
	workspacesInfrastructure "backend/internal/workspaces/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
	repo := infrastructure.NewAssignmentRepository()
	repo.AutoMigrate()
	userRepo := usersInfrastructure.NewUserRepository()
	workspaceRepo := workspacesInfrastructure.NewWorkspaceRepository()

	createAssignment := application.NewCreateAssignment(repo).WithRepositories(userRepo, workspaceRepo)
	getAssignmentByID := application.NewGetAssignmentByID(repo)
	listAssignmentsByUserID := application.NewListAssignmentsByUserID(repo)
	updateAssignment := application.NewUpdateAssignment(repo)

	handler := NewAssignmentHandler(createAssignment, getAssignmentByID, listAssignmentsByUserID, updateAssignment)

	assignments := r.Group("/assignments")
	{
		assignments.POST("", handler.CreateAssignment)
		assignments.GET("/:id", handler.GetAssignmentByID)
		assignments.GET("", handler.ListAssignmentsByUserID)
		//nolint:godox // TODO RF04: Proteger esta ruta para que solo admin pueda actualizar vinculaciones cuando auth este integrado.
		assignments.PUT("/:id/admin", handler.UpdateAssignment)
	}
}
