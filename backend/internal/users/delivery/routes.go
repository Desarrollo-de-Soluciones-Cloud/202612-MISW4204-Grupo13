package delivery

import (
	"backend/internal/users/application"
	"backend/internal/users/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
	repo := infrastructure.NewUserRepository()
	repo.AutoMigrate()
	createUser := application.NewCreateUser(repo)
	listUsers := application.NewListUsers(repo)
	listUsersByRole := application.NewListUsersByRole(repo)
	getUserByID := application.NewGetUserByID(repo)
	updateUser := application.NewUpdateUser(repo)
	changeUserRole := application.NewChangeUserRole(repo)
	handler := NewUserHandler(createUser, listUsers, listUsersByRole, getUserByID, updateUser, changeUserRole)
	users := r.Group("/users")
	{
		users.POST("", handler.CreateUser)
		users.GET("", handler.ListUsers)
		users.GET("/:id", handler.GetUserByID)
		users.PUT("/:id", handler.UpdateUser)
		users.PATCH("/:id/role", handler.ChangeUserRole)
	}
}
