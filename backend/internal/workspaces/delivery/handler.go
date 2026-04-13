package delivery

import (
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/workspaces/application"
	"backend/internal/workspaces/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	createWorkspace         *application.CreateWorkspace
	listWorkspaces          *application.ListWorkspaces
	listWorkspacesByPeriod  *application.ListWorkspacesByPeriod
	getWorkspaceByID        *application.GetWorkspaceByID
	updateWorkspace         *application.UpdateWorkspace
	deleteWorkspace         *application.DeleteWorkspace
}

func NewWorkspaceHandler(
	createWorkspace *application.CreateWorkspace,
	listWorkspaces *application.ListWorkspaces,
	listWorkspacesByPeriod *application.ListWorkspacesByPeriod,
	getWorkspaceByID *application.GetWorkspaceByID,
	updateWorkspace *application.UpdateWorkspace,
	deleteWorkspace *application.DeleteWorkspace,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		createWorkspace:        createWorkspace,
		listWorkspaces:         listWorkspaces,
		listWorkspacesByPeriod: listWorkspacesByPeriod,
		getWorkspaceByID:       getWorkspaceByID,
		updateWorkspace:        updateWorkspace,
		deleteWorkspace:        deleteWorkspace,
	}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	var req CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	input := application.CreateWorkspaceInput{
		PeriodID:     req.PeriodID,
		UserID:       req.UserID,
		Name:         req.Name,
		Type:         req.Type,
		InitialDate:  req.InitialDate,
		FinalDate:    req.FinalDate,
		Observations: req.Observations,
		State:        req.State,
	}

	output, err := h.createWorkspace.Execute(input)
	if err != nil {
		switch {
		case isWorkspaceValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, CreateWorkspaceResponse{
		ID:           output.ID,
		PeriodID:     output.PeriodID,
		UserID:       output.UserID,
		Name:         output.Name,
		Type:         output.Type,
		InitialDate:  output.InitialDate,
		FinalDate:    output.FinalDate,
		Observations: output.Observations,
		State:        output.State,
	})
}

func (h *WorkspaceHandler) ListWorkspaces(c *gin.Context) {
	periodID := c.Query("period_id")
	if periodID != "" {
		h.listWorkspacesByPeriodHandler(c, periodID)
		return
	}

	output, err := h.listWorkspaces.Execute()
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	workspaces := make([]WorkspaceResponse, len(output.Workspaces))
	for i, w := range output.Workspaces {
		workspaces[i] = WorkspaceResponse{
			ID:           w.ID,
			PeriodID:     w.PeriodID,
			UserID:       w.UserID,
			Name:         w.Name,
			Type:         w.Type,
			InitialDate:  w.InitialDate,
			FinalDate:    w.FinalDate,
			Observations: w.Observations,
			State:        w.State,
		}
	}

	c.JSON(http.StatusOK, ListWorkspacesResponse{Workspaces: workspaces})
}

func (h *WorkspaceHandler) listWorkspacesByPeriodHandler(c *gin.Context, rawPeriodID string) {
	periodID, err := sharedHelpers.ParseResourceID(rawPeriodID)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.listWorkspacesByPeriod.Execute(application.ListWorkspacesByPeriodInput{
		PeriodID: periodID,
	})
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	workspaces := make([]WorkspaceResponse, len(output.Workspaces))
	for i, w := range output.Workspaces {
		workspaces[i] = WorkspaceResponse{
			ID:           w.ID,
			PeriodID:     w.PeriodID,
			UserID:       w.UserID,
			Name:         w.Name,
			Type:         w.Type,
			InitialDate:  w.InitialDate,
			FinalDate:    w.FinalDate,
			Observations: w.Observations,
			State:        w.State,
		}
	}

	c.JSON(http.StatusOK, ListWorkspacesResponse{Workspaces: workspaces})
}

func (h *WorkspaceHandler) GetWorkspaceByID(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.getWorkspaceByID.Execute(application.GetWorkspaceByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWorkspaceNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, WorkspaceResponse{
		ID:           output.ID,
		PeriodID:     output.PeriodID,
		UserID:       output.UserID,
		Name:         output.Name,
		Type:         output.Type,
		InitialDate:  output.InitialDate,
		FinalDate:    output.FinalDate,
		Observations: output.Observations,
		State:        output.State,
	})
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	var req UpdateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	output, err := h.updateWorkspace.Execute(application.UpdateWorkspaceInput{
		ID:           id,
		PeriodID:     req.PeriodID,
		Name:         req.Name,
		Type:         req.Type,
		InitialDate:  req.InitialDate,
		FinalDate:    req.FinalDate,
		Observations: req.Observations,
		State:        req.State,
	})
	if err != nil {
		switch {
		case errors.Is(err, errors.New("workspace not found")):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case isWorkspaceValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, WorkspaceResponse{
		ID:           output.ID,
		PeriodID:     output.PeriodID,
		UserID:       output.UserID,
		Name:         output.Name,
		Type:         output.Type,
		InitialDate:  output.InitialDate,
		FinalDate:    output.FinalDate,
		Observations: output.Observations,
		State:        output.State,
	})
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	err = h.deleteWorkspace.Execute(application.DeleteWorkspaceInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, errors.New("workspace not found")):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
