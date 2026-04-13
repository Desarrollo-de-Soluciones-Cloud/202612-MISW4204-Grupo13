package delivery

import (
	"backend/internal/periods/application"
	"backend/internal/periods/infrastructure"
	weeksApplication "backend/internal/weeks/application"
	weeksInfrastructure "backend/internal/weeks/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
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
		periods.POST("", handler.CreatePeriod)
		periods.GET("", handler.ListPeriods)
		periods.GET("/:id", handler.GetPeriodByID)
		periods.PATCH("/:id", handler.UpdatePeriod)
		periods.PATCH("/:id/close", handler.ClosePeriod)
	}
}
// merge a develop
