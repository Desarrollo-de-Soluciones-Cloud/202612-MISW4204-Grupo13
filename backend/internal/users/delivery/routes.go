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
    handler := NewUserHandler(createUser, listUsers)
    users := r.Group("/users")
    {
        users.POST("", handler.CreateUser)
        users.GET("", handler.ListUsers)
    }
}