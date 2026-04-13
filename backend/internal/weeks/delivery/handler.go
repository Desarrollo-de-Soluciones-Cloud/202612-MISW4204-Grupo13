package delivery

import (
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/weeks/application"
	"backend/internal/weeks/domain"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WeekHandler struct {
	listWeeksByPeriod         *application.ListWeeksByPeriod
	getWeekByPeriodAndNumber  *application.GetWeekByPeriodAndNumber
}

func NewWeekHandler(
	listWeeksByPeriod *application.ListWeeksByPeriod,
	getWeekByPeriodAndNumber *application.GetWeekByPeriodAndNumber,
) *WeekHandler {
	return &WeekHandler{
		listWeeksByPeriod:        listWeeksByPeriod,
		getWeekByPeriodAndNumber: getWeekByPeriodAndNumber,
	}
}

func (h *WeekHandler) ListWeeksByPeriod(c *gin.Context) {
	periodID, err := sharedHelpers.ParseResourceID(c.Param("periodId"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.listWeeksByPeriod.Execute(application.ListWeeksByPeriodInput{
		PeriodID: periodID,
	})
	if err != nil {
		switch {
		case isWeekValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	weeks := make([]WeekResponse, len(output.Weeks))
	for i, week := range output.Weeks {
		weeks[i] = WeekResponse{
			ID:          week.ID,
			PeriodID:    week.PeriodID,
			Number:      week.Number,
			InitialDate: week.InitialDate,
			FinalDate:   week.FinalDate,
		}
	}

	c.JSON(http.StatusOK, ListWeeksResponse{Weeks: weeks})
}

func (h *WeekHandler) GetWeekByPeriodAndNumber(c *gin.Context) {
	periodID, err := sharedHelpers.ParseResourceID(c.Param("periodId"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	number, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, domain.ErrWeekNumberInvalid)
		return
	}

	output, err := h.getWeekByPeriodAndNumber.Execute(application.GetWeekByPeriodAndNumberInput{
		PeriodID: periodID,
		Number:   number,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWeekNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case isWeekValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, WeekResponse{
		ID:          output.ID,
		PeriodID:    output.PeriodID,
		Number:      output.Number,
		InitialDate: output.InitialDate,
		FinalDate:   output.FinalDate,
	})
}

func isWeekValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrWeekPeriodIDRequired) ||
		errors.Is(err, domain.ErrWeekNumberInvalid) ||
		errors.Is(err, domain.ErrWeekInitialDateRequired) ||
		errors.Is(err, domain.ErrWeekInitialDateWrongFormat) ||
		errors.Is(err, domain.ErrWeekInitialDateMustBeMonday) ||
		errors.Is(err, domain.ErrWeekFinalDateRequired) ||
		errors.Is(err, domain.ErrWeekFinalDateWrongFormat) ||
		errors.Is(err, domain.ErrWeekFinalDateMustBeSunday) ||
		errors.Is(err, domain.ErrWeekDateRangeInvalid) ||
		errors.Is(err, domain.ErrWeekFinalDateMismatch) ||
		errors.Is(err, domain.ErrWeekCountInvalid)
}
