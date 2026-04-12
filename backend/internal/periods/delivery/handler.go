package delivery

import (
	"backend/internal/periods/application"
	"backend/internal/periods/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PeriodHandler struct {
	createPeriod       *application.CreatePeriod
	listPeriods        *application.ListPeriods
	listPeriodsByState *application.ListPeriodsByState
	getPeriodByID      *application.GetPeriodByID
	updatePeriod       *application.UpdatePeriod
	deletePeriod       *application.DeletePeriod
}

func NewPeriodHandler(
	createPeriod *application.CreatePeriod,
	listPeriods *application.ListPeriods,
	listPeriodsByState *application.ListPeriodsByState,
	getPeriodByID *application.GetPeriodByID,
	updatePeriod *application.UpdatePeriod,
	deletePeriod *application.DeletePeriod,
) *PeriodHandler {
	return &PeriodHandler{
		createPeriod:       createPeriod,
		listPeriods:        listPeriods,
		listPeriodsByState: listPeriodsByState,
		getPeriodByID:      getPeriodByID,
		updatePeriod:       updatePeriod,
		deletePeriod:       deletePeriod,
	}
}

func (h *PeriodHandler) CreatePeriod(c *gin.Context) {
	var req CreatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	input := application.CreatePeriodInput{
		Name:                 req.Name,
		InitialDate:          req.InitialDate,
		FinalDate:            req.FinalDate,
		InscriptionFinalDate: req.InscriptionFinalDate,
		PeriodState:          domain.PeriodState(req.PeriodState),
	}

	output, err := h.createPeriod.Execute(input)
	if err != nil {
		switch {
		case isPeriodValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, CreatePeriodResponse{
		ID:                   output.ID,
		Name:                 output.Name,
		InitialDate:          output.InitialDate,
		FinalDate:            output.FinalDate,
		InscriptionFinalDate: output.InscriptionFinalDate,
		PeriodState:          string(output.PeriodState),
	})
}

func (h *PeriodHandler) ListPeriods(c *gin.Context) {
	stateFilter := c.Query("state")
	if stateFilter != "" {
		h.listPeriodsByStateHandler(c, stateFilter)
		return
	}

	output, err := h.listPeriods.Execute()
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	periods := make([]PeriodResponse, len(output.Periods))
	for i, p := range output.Periods {
		periods[i] = PeriodResponse{
			ID:                   p.ID,
			Name:                 p.Name,
			InitialDate:          p.InitialDate,
			FinalDate:            p.FinalDate,
			InscriptionFinalDate: p.InscriptionFinalDate,
			PeriodState:          string(p.PeriodState),
		}
	}

	c.JSON(http.StatusOK, ListPeriodsResponse{Periods: periods})
}

func (h *PeriodHandler) listPeriodsByStateHandler(c *gin.Context, rawState string) {
	output, err := h.listPeriodsByState.Execute(application.ListPeriodsByStateInput{
		PeriodState: domain.PeriodState(rawState),
	})
	if err != nil {
		switch {
		case isPeriodValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	periods := make([]PeriodResponse, len(output.Periods))
	for i, p := range output.Periods {
		periods[i] = PeriodResponse{
			ID:                   p.ID,
			Name:                 p.Name,
			InitialDate:          p.InitialDate,
			FinalDate:            p.FinalDate,
			InscriptionFinalDate: p.InscriptionFinalDate,
			PeriodState:          string(p.PeriodState),
		}
	}

	c.JSON(http.StatusOK, ListPeriodsResponse{Periods: periods})
}

func (h *PeriodHandler) GetPeriodByID(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.getPeriodByID.Execute(application.GetPeriodByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPeriodNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, PeriodResponse{
		ID:                   output.ID,
		Name:                 output.Name,
		InitialDate:          output.InitialDate,
		FinalDate:            output.FinalDate,
		InscriptionFinalDate: output.InscriptionFinalDate,
		PeriodState:          string(output.PeriodState),
	})
}

func (h *PeriodHandler) UpdatePeriod(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	var req UpdatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	output, err := h.updatePeriod.Execute(application.UpdatePeriodInput{
		ID:                   id,
		Name:                 req.Name,
		InitialDate:          req.InitialDate,
		FinalDate:            req.FinalDate,
		InscriptionFinalDate: req.InscriptionFinalDate,
		PeriodState:          domain.PeriodState(req.PeriodState),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPeriodNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case isPeriodValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, PeriodResponse{
		ID:                   output.ID,
		Name:                 output.Name,
		InitialDate:          output.InitialDate,
		FinalDate:            output.FinalDate,
		InscriptionFinalDate: output.InscriptionFinalDate,
		PeriodState:          string(output.PeriodState),
	})
}

func (h *PeriodHandler) DeletePeriod(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	err = h.deletePeriod.Execute(application.DeletePeriodInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPeriodNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func isPeriodValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrPeriodNameRequired) ||
		errors.Is(err, domain.ErrPeriodNameWrongFormat) ||
		errors.Is(err, domain.ErrPeriodInitialDateRequired) ||
		errors.Is(err, domain.ErrPeriodInitialDateInvalid) ||
		errors.Is(err, domain.ErrPeriodFinalDateRequired) ||
		errors.Is(err, domain.ErrPeriodFinalDateInvalid) ||
		errors.Is(err, domain.ErrPeriodInscriptionFinalDateRequired) ||
		errors.Is(err, domain.ErrPeriodInscriptionFinalDateInvalid) ||
		errors.Is(err, domain.ErrPeriodDateSequenceInvalid) ||
		errors.Is(err, domain.ErrPeriodStateRequired) ||
		errors.Is(err, domain.ErrPeriodStateInvalid)
}
