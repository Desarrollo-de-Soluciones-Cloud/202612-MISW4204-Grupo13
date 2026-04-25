package delivery

import (
	"backend/internal/users/application"
	"backend/internal/users/domain"
	"backend/internal/users/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...domain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer) {
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
		users.Use(authorizer.RequireAuthentication())

		adminUsers := users.Group("")
		adminUsers.Use(authorizer.RequireRoles(domain.RoleAdmin))
		adminUsers.POST("", handler.CreateUser)
		adminUsers.PUT("/:id", handler.UpdateUser)
		adminUsers.PATCH("/:id/role", handler.ChangeUserRole)
		adminUsers.GET("", handler.ListUsers)
	}
}
