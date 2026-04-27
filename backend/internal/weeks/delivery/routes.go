package delivery

import (
	usersDomain "backend/internal/users/domain"
	"backend/internal/weeks/application"
	"backend/internal/weeks/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...usersDomain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer) {
	repo := infrastructure.NewWeekRepository()
	repo.AutoMigrate()
	listWeeksByPeriod := application.NewListWeeksByPeriod(repo)
	getWeekByPeriodAndNumber := application.NewGetWeekByPeriodAndNumber(repo)
	handler := NewWeekHandler(listWeeksByPeriod, getWeekByPeriodAndNumber)

	weeks := r.Group("/weeks")
	{
		weeks.Use(authorizer.RequireAuthentication())

		adminAndProfessorWeeks := weeks.Group("")
		adminAndProfessorWeeks.Use(authorizer.RequireRoles(usersDomain.RoleAdmin, usersDomain.RoleProfessor))
		adminAndProfessorWeeks.GET("/periods/:periodId", handler.ListWeeksByPeriod)
		adminAndProfessorWeeks.GET("/:number/periods/:periodId", handler.GetWeekByPeriodAndNumber)
	}
}
