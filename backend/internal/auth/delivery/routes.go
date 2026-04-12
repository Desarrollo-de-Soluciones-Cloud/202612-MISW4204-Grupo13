package delivery

import (
	"backend/internal/auth/application"
	"backend/internal/auth/infrastructure"
	"backend/internal/shared/config"
	usersApplication "backend/internal/users/application"
	usersInfrastructure "backend/internal/users/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter, cfg *config.Config) *AuthHandler {
	userRepository := usersInfrastructure.NewUserRepository()
	getUserByEmail := usersApplication.NewGetUserByEmail(userRepository)
	userReader := infrastructure.NewUserReader(getUserByEmail)
	tokenManager := infrastructure.NewTokenManager(cfg.JWTSecret, cfg.JWTExpirationMinutes)
	signIn := application.NewSignIn(userReader, tokenManager)
	validateToken := application.NewValidateToken(tokenManager)
	handler := NewAuthHandler(signIn, validateToken)

	auth := r.Group("/auth")
	{
		auth.POST("/sign-in", handler.SignIn)

		authenticated := auth.Group("")
		authenticated.Use(handler.RequireAuthentication())
		authenticated.GET("/me", handler.GetCurrentUser)
	}

	return handler
}
