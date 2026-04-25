package delivery

import (
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	usersDomain "backend/internal/users/domain"
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
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	var req CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	switch currentUser.GlobalRole {
	case usersDomain.RoleAdmin:
		// Admin can create workspaces for any professor.
	case usersDomain.RoleProfessor:
		if req.UserID != currentUser.ID {
			sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
			return
		}
	default:
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
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
		case errors.Is(err, domain.ErrWorkspacePeriodNotFound), errors.Is(err, domain.ErrWorkspaceUserNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrWorkspacePeriodClosed), errors.Is(err, domain.ErrWorkspaceInscriptionClosed):
			sharedHelpers.RespondWithError(c, http.StatusConflict, err)
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
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	periodID := c.Query("period_id")
	if periodID != "" {
		h.listWorkspacesByPeriodHandler(c, currentUser, periodID)
		return
	}

	output, err := h.listWorkspaces.Execute()
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	workspaces := make([]WorkspaceResponse, 0, len(output.Workspaces))
	for _, w := range output.Workspaces {
		if !canAccessWorkspace(currentUser.GlobalRole, currentUser.ID, w.UserID) {
			continue
		}

		workspaces = append(workspaces, WorkspaceResponse{
			ID:           w.ID,
			PeriodID:     w.PeriodID,
			UserID:       w.UserID,
			Name:         w.Name,
			Type:         w.Type,
			InitialDate:  w.InitialDate,
			FinalDate:    w.FinalDate,
			Observations: w.Observations,
			State:        w.State,
		})
	}

	c.JSON(http.StatusOK, ListWorkspacesResponse{Workspaces: workspaces})
}

func (h *WorkspaceHandler) listWorkspacesByPeriodHandler(c *gin.Context, currentUser authDomain.AuthenticatedUser, rawPeriodID string) {
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

	workspaces := make([]WorkspaceResponse, 0, len(output.Workspaces))
	for _, w := range output.Workspaces {
		if !canAccessWorkspace(currentUser.GlobalRole, currentUser.ID, w.UserID) {
			continue
		}

		workspaces = append(workspaces, WorkspaceResponse{
			ID:           w.ID,
			PeriodID:     w.PeriodID,
			UserID:       w.UserID,
			Name:         w.Name,
			Type:         w.Type,
			InitialDate:  w.InitialDate,
			FinalDate:    w.FinalDate,
			Observations: w.Observations,
			State:        w.State,
		})
	}

	c.JSON(http.StatusOK, ListWorkspacesResponse{Workspaces: workspaces})
}

func (h *WorkspaceHandler) GetWorkspaceByID(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

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
	if !canAccessWorkspace(currentUser.GlobalRole, currentUser.ID, output.UserID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
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
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	existing, err := h.getWorkspaceByID.Execute(application.GetWorkspaceByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWorkspaceNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}
	if !canAccessWorkspace(currentUser.GlobalRole, currentUser.ID, existing.UserID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
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
		case errors.Is(err, domain.ErrWorkspaceNotFound), errors.Is(err, domain.ErrWorkspacePeriodNotFound), errors.Is(err, domain.ErrWorkspaceUserNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrWorkspaceClosedUpdateForbidden):
			sharedHelpers.RespondWithError(c, http.StatusConflict, err)
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
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	existing, err := h.getWorkspaceByID.Execute(application.GetWorkspaceByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWorkspaceNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}
	if !canAccessWorkspace(currentUser.GlobalRole, currentUser.ID, existing.UserID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	err = h.deleteWorkspace.Execute(application.DeleteWorkspaceInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWorkspaceNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func canAccessWorkspace(role usersDomain.UserRole, currentUserID, workspaceUserID uint) bool {
	switch role {
	case usersDomain.RoleAdmin:
		return true
	case usersDomain.RoleProfessor:
		return currentUserID == workspaceUserID
	default:
		return false
	}
}
