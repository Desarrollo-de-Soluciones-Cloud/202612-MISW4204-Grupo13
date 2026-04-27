package delivery

import (
	"backend/internal/periods/application"
	"backend/internal/periods/infrastructure"
	usersDomain "backend/internal/users/domain"
	weeksApplication "backend/internal/weeks/application"
	weeksInfrastructure "backend/internal/weeks/infrastructure"

	"github.com/gin-gonic/gin"
)

type RouteAuthorizer interface {
	RequireAuthentication() gin.HandlerFunc
	RequireRoles(...usersDomain.UserRole) gin.HandlerFunc
}

func SetupRoutes(r gin.IRouter, authorizer RouteAuthorizer) {
	repo := infrastructure.NewPeriodRepository()
	repo.AutoMigrate()
	weeksRepo := weeksInfrastructure.NewWeekRepository()
	weeksRepo.AutoMigrate()
	createWeeksForPeriod := weeksApplication.NewCreateWeeksForPeriod(weeksRepo)
	createPeriod := application.NewCreatePeriod(repo, createWeeksForPeriod)
	listPeriods := application.NewListPeriods(repo)
	listPeriodsByState := application.NewListPeriodsByState(repo)
	getPeriodByID := application.NewGetPeriodByID(repo)
	updatePeriod := application.NewUpdatePeriod(repo)
	closePeriod := application.NewClosePeriod(repo)
	handler := NewPeriodHandler(createPeriod, listPeriods, listPeriodsByState, getPeriodByID, updatePeriod, closePeriod)
	periods := r.Group("/periods")
	{
		periods.Use(authorizer.RequireAuthentication())

		adminPeriods := periods.Group("")
		adminPeriods.Use(authorizer.RequireRoles(usersDomain.RoleAdmin))
		adminPeriods.POST("", handler.CreatePeriod)
		adminPeriods.GET("", handler.ListPeriods)
		adminPeriods.PATCH("/:id", handler.UpdatePeriod)
		adminPeriods.PATCH("/:id/close", handler.ClosePeriod)
	}
}
// merge a develop
