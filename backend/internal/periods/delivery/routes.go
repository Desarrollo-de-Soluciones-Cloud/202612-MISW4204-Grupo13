package delivery

import (
	"backend/internal/periods/application"
	"backend/internal/periods/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r gin.IRouter) {
	repo := infrastructure.NewPeriodRepository()
	repo.AutoMigrate()
	createPeriod := application.NewCreatePeriod(repo)
	listPeriods := application.NewListPeriods(repo)
	listPeriodsByState := application.NewListPeriodsByState(repo)
	getPeriodByID := application.NewGetPeriodByID(repo)
	updatePeriod := application.NewUpdatePeriod(repo)
	deletePeriod := application.NewDeletePeriod(repo)
	handler := NewPeriodHandler(createPeriod, listPeriods, listPeriodsByState, getPeriodByID, updatePeriod, deletePeriod)
	periods := r.Group("/periods")
	{
		periods.POST("", handler.CreatePeriod)
		periods.GET("", handler.ListPeriods)
		periods.GET("/:id", handler.GetPeriodByID)
		periods.PUT("/:id", handler.UpdatePeriod)
		periods.DELETE("/:id", handler.DeletePeriod)
	}
}
// merge into develop