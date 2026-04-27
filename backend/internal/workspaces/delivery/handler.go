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
	createWorkspace                          *application.CreateWorkspace
	listWorkspaces                           *application.ListWorkspaces
	listWorkspacesByPeriod                   *application.ListWorkspacesByPeriod
	getWorkspaceByID                         *application.GetWorkspaceByID
	updateWorkspace                          *application.UpdateWorkspace
	deleteWorkspace                          *application.DeleteWorkspace
	closeWorkspace                           *application.CloseWorkspace
	listWorkspaceMonitorsAndAssistants       *application.ListWorkspaceMonitorsAndAssistants
}

func NewWorkspaceHandler(
	createWorkspace *application.CreateWorkspace,
	listWorkspaces *application.ListWorkspaces,
	listWorkspacesByPeriod *application.ListWorkspacesByPeriod,
	getWorkspaceByID *application.GetWorkspaceByID,
	updateWorkspace *application.UpdateWorkspace,
	deleteWorkspace *application.DeleteWorkspace,
	closeWorkspace *application.CloseWorkspace,
	listWorkspaceMonitorsAndAssistants *application.ListWorkspaceMonitorsAndAssistants,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		createWorkspace:                          createWorkspace,
		listWorkspaces:                           listWorkspaces,
		listWorkspacesByPeriod:                   listWorkspacesByPeriod,
		getWorkspaceByID:                         getWorkspaceByID,
		updateWorkspace:                          updateWorkspace,
		deleteWorkspace:                          deleteWorkspace,
		closeWorkspace:                           closeWorkspace,
		listWorkspaceMonitorsAndAssistants:       listWorkspaceMonitorsAndAssistants,
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
		case errors.Is(err, domain.ErrWorkspacePeriodClosed), errors.Is(err, domain.ErrWorkspaceInscriptionClosed), errors.Is(err, domain.ErrWorkspaceInitialDateOutOfRange), errors.Is(err, domain.ErrWorkspaceFinalDateOutOfRange):
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

	// Check for user_id filter
	userIDQuery := c.Query("user_id")
	if userIDQuery != "" {
		// Professors cannot filter by a different user_id
		if currentUser.GlobalRole == usersDomain.RoleProfessor {
			requestedUserID, err := sharedHelpers.ParseResourceID(userIDQuery)
			if err != nil {
				sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
				return
			}
			if requestedUserID != currentUser.ID {
				sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
				return
			}
		}
		// Admin can filter by any user_id - this will be handled in filtering below
	}

	// Get type and state filters
	typeFilter := c.Query("type")
	stateFilter := c.Query("state")

	periodID := c.Query("period_id")
	if periodID != "" {
		h.listWorkspacesByPeriodHandler(c, currentUser, periodID, userIDQuery, typeFilter, stateFilter)
		return
	}

	output, err := h.listWorkspaces.Execute()
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	var requestedUserID uint
	if userIDQuery != "" {
		var err error
		requestedUserID, err = sharedHelpers.ParseResourceID(userIDQuery)
		if err != nil {
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
	}

	workspaces := make([]WorkspaceResponse, 0, len(output.Workspaces))
	for _, w := range output.Workspaces {
		if !canAccessWorkspace(currentUser.GlobalRole, currentUser.ID, w.UserID) {
			continue
		}

		// Apply user_id filter if provided
		if requestedUserID > 0 && w.UserID != requestedUserID {
			continue
		}

		// For professors without explicit user_id filter, only show their own workspaces
		if currentUser.GlobalRole == usersDomain.RoleProfessor && requestedUserID == 0 {
			if w.UserID != currentUser.ID {
				continue
			}
		}

		// Apply type filter if provided
		if typeFilter != "" && w.Type != typeFilter {
			continue
		}

		// Apply state filter if provided
		if stateFilter != "" && w.State != stateFilter {
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

func (h *WorkspaceHandler) listWorkspacesByPeriodHandler(c *gin.Context, currentUser authDomain.AuthenticatedUser, rawPeriodID string, userIDQuery string, typeFilter string, stateFilter string) {
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

	var requestedUserID uint
	if userIDQuery != "" {
		var err error
		requestedUserID, err = sharedHelpers.ParseResourceID(userIDQuery)
		if err != nil {
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
			return
		}
	}

	workspaces := make([]WorkspaceResponse, 0, len(output.Workspaces))
	for _, w := range output.Workspaces {
		if !canAccessWorkspace(currentUser.GlobalRole, currentUser.ID, w.UserID) {
			continue
		}

		// Apply user_id filter if provided
		if requestedUserID > 0 && w.UserID != requestedUserID {
			continue
		}

		// For professors without explicit user_id filter, only show their own workspaces
		if currentUser.GlobalRole == usersDomain.RoleProfessor && requestedUserID == 0 {
			if w.UserID != currentUser.ID {
				continue
			}
		}

		// Apply type filter if provided
		if typeFilter != "" && w.Type != typeFilter {
			continue
		}

		// Apply state filter if provided
		if stateFilter != "" && w.State != stateFilter {
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
		UserID:       req.UserID,
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
		case errors.Is(err, domain.ErrWorkspaceClosedUpdateForbidden), errors.Is(err, domain.ErrWorkspacePeriodClosed), errors.Is(err, domain.ErrWorkspaceInitialDateOutOfRange), errors.Is(err, domain.ErrWorkspaceFinalDateOutOfRange), errors.Is(err, domain.ErrWorkspaceUserIDChangeNotAllowed):
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

func (h *WorkspaceHandler) CloseWorkspace(c *gin.Context) {
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

	// Get workspace to verify access
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

	output, err := h.closeWorkspace.Execute(application.CloseWorkspaceInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWorkspaceNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrWorkspaceUserNotFound), errors.Is(err, domain.ErrWorkspaceUserNotProfessor):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, CloseWorkspaceResponse{
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

func (h *WorkspaceHandler) ListWorkspaceMonitorsAndAssistants(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	if currentUser.GlobalRole != usersDomain.RoleProfessor {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	output, err := h.listWorkspaceMonitorsAndAssistants.Execute(application.ListWorkspaceMonitorsAndAssistantsInput{
		ProfessorID: currentUser.ID,
	})
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	monitors := make([]MonitorAssistantResponse, len(output.Monitors))
	for i, m := range output.Monitors {
		monitors[i] = MonitorAssistantResponse{
			ID:          m.ID,
			Name:        m.Name,
			Email:       m.Email,
			Role:        m.Role,
			WeeklyHours: m.WeeklyHours,
		}
	}

	assistants := make([]MonitorAssistantResponse, len(output.Assistants))
	for i, a := range output.Assistants {
		assistants[i] = MonitorAssistantResponse{
			ID:          a.ID,
			Name:        a.Name,
			Email:       a.Email,
			Role:        a.Role,
			WeeklyHours: a.WeeklyHours,
		}
	}

	c.JSON(http.StatusOK, ListWorkspaceMonitorsAndAssistantsResponse{
		Monitors:   monitors,
		Assistants: assistants,
	})
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
