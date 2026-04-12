package delivery

import (
	"backend/internal/assignments/application"
	"backend/internal/assignments/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
	repo := infrastructure.NewAssignmentRepository()
	repo.AutoMigrate()

	createAssignment := application.NewCreateAssignment(repo)
	getAssignmentByID := application.NewGetAssignmentByID(repo)
	listAssignmentsByUserID := application.NewListAssignmentsByUserID(repo)

	handler := NewAssignmentHandler(createAssignment, getAssignmentByID, listAssignmentsByUserID)

	assignments := r.Group("/assignments")
	{
		assignments.POST("", handler.CreateAssignment)
		assignments.GET("/:id", handler.GetAssignmentByID)
		assignments.GET("", handler.ListAssignmentsByUserID)
	}
}
