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

	weeks := r.Group("/weeks")
	{
		weeks.GET("/periods/:periodId", handler.ListWeeksByPeriod)
		weeks.GET("/:number/periods/:periodId", handler.GetWeekByPeriodAndNumber)
	}
}
