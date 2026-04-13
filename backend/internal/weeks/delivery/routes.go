package delivery

import (
	"backend/internal/weeks/application"
	"backend/internal/weeks/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
	repo := infrastructure.NewWeekRepository()
	repo.AutoMigrate()
	listWeeksByPeriod := application.NewListWeeksByPeriod(repo)
	getWeekByPeriodAndNumber := application.NewGetWeekByPeriodAndNumber(repo)
	handler := NewWeekHandler(listWeeksByPeriod, getWeekByPeriodAndNumber)

	periods := r.Group("/periods/:periodId/weeks")
	{
		periods.GET("", handler.ListWeeksByPeriod)
		periods.GET("/:number", handler.GetWeekByPeriodAndNumber)
	}
}
