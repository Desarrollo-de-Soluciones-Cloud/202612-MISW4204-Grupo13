package delivery

import (
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssignmentHandler struct {
	createAssignment        *application.CreateAssignment
	getAssignmentByID       *application.GetAssignmentByID
	listAssignmentsByUserID *application.ListAssignmentsByUserID
	updateAssignment        *application.UpdateAssignment
	workspaceReader         AssignmentWorkspaceReader
}

type AssignmentWorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

func NewAssignmentHandler(
	createAssignment *application.CreateAssignment,
	getAssignmentByID *application.GetAssignmentByID,
	listAssignmentsByUserID *application.ListAssignmentsByUserID,
	updateAssignment *application.UpdateAssignment,
	workspaceReader AssignmentWorkspaceReader,
) *AssignmentHandler {
	return &AssignmentHandler{
		createAssignment:        createAssignment,
		getAssignmentByID:       getAssignmentByID,
		listAssignmentsByUserID: listAssignmentsByUserID,
		updateAssignment:        updateAssignment,
		workspaceReader:         workspaceReader,
	}
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	switch currentUser.GlobalRole {
	case usersDomain.RoleAdmin:
		// Admin can create assignments on any workspace.
	case usersDomain.RoleProfessor:
		workspace, err := h.workspaceReader.FindByID(req.WorkspaceID)
		if err != nil {
			switch {
			case isWorkspaceNotFoundError(err):
				sharedHelpers.RespondWithError(c, http.StatusNotFound, domain.ErrAssignmentWorkspaceNotFound)
			default:
				sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
			}
			return
		}

		if workspace.UserID != currentUser.ID {
			sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
			return
		}
	default:
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	output, err := h.createAssignment.Execute(application.CreateAssignmentInput{
		UserID:      req.UserID,
		WorkspaceID: req.WorkspaceID,
		Role:        domain.AssignmentRole(req.Role),
		WeeklyHours: req.WeeklyHours,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAssignmentUserNotFound),
			errors.Is(err, domain.ErrAssignmentWorkspaceNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrAssignmentWorkspaceClosed),
			errors.Is(err, domain.ErrAssignmentAlreadyExists):
			sharedHelpers.RespondWithError(c, http.StatusConflict, err)
		case isAssignmentValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, AssignmentResponse{
		ID:          output.ID,
		UserID:      output.UserID,
		WorkspaceID: output.WorkspaceID,
		Role:        string(output.Role),
		WeeklyHours: output.WeeklyHours,
	})
}

func (h *AssignmentHandler) GetAssignmentByID(c *gin.Context) {
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

	output, err := h.getAssignmentByID.Execute(application.GetAssignmentByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	if !h.canReadAssignment(currentUser.GlobalRole, currentUser.ID, output.UserID, output.WorkspaceID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	c.JSON(http.StatusOK, AssignmentResponse{
		ID:          output.ID,
		UserID:      output.UserID,
		WorkspaceID: output.WorkspaceID,
		Role:        string(output.Role),
		WeeklyHours: output.WeeklyHours,
	})
}

func (h *AssignmentHandler) UpdateAssignment(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	var req UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	output, err := h.updateAssignment.Execute(application.UpdateAssignmentInput{
		ID:          id,
		Role:        domain.AssignmentRole(req.Role),
		WeeklyHours: req.WeeklyHours,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrAssignmentAlreadyExists):
			sharedHelpers.RespondWithError(c, http.StatusConflict, err)
		case isAssignmentValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, AssignmentResponse{
		ID:          output.ID,
		UserID:      output.UserID,
		WorkspaceID: output.WorkspaceID,
		Role:        string(output.Role),
		WeeklyHours: output.WeeklyHours,
	})
}

func (h *AssignmentHandler) ListAssignmentsByUserID(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	rawUserID := c.Query("user_id")
	parsedID, err := strconv.ParseUint(rawUserID, 10, 32)
	if err != nil || parsedID == 0 {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, domain.ErrAssignmentUserIDRequired)
		return
	}

	targetUserID := uint(parsedID)
	if (currentUser.GlobalRole == usersDomain.RoleMonitor || currentUser.GlobalRole == usersDomain.RoleAssistant) && currentUser.ID != targetUserID {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	output, err := h.listAssignmentsByUserID.Execute(application.ListAssignmentsByUserIDInput{
		UserID: targetUserID,
	})
	if err != nil {
		switch {
		case isAssignmentValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	assignments := make([]AssignmentResponse, 0, len(output.Assignments))
	for _, a := range output.Assignments {
		if !h.canReadAssignment(currentUser.GlobalRole, currentUser.ID, a.UserID, a.WorkspaceID) {
			sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
			return
		}

		assignments = append(assignments, AssignmentResponse{
			ID:          a.ID,
			UserID:      a.UserID,
			WorkspaceID: a.WorkspaceID,
			Role:        string(a.Role),
			WeeklyHours: a.WeeklyHours,
		})
	}

	c.JSON(http.StatusOK, ListAssignmentsResponse{Assignments: assignments})
}

func (h *AssignmentHandler) canReadAssignment(role usersDomain.UserRole, currentUserID, assignmentUserID, workspaceID uint) bool {
	switch role {
	case usersDomain.RoleAdmin:
		return true
	case usersDomain.RoleMonitor, usersDomain.RoleAssistant:
		return currentUserID == assignmentUserID
	case usersDomain.RoleProfessor:
		workspace, err := h.workspaceReader.FindByID(workspaceID)
		if err != nil {
			return false
		}
		return workspace.UserID == currentUserID
	default:
		return false
	}
}

func isWorkspaceNotFoundError(err error) bool {
	return errors.Is(err, workspacesDomain.ErrWorkspaceNotFound) || err.Error() == workspacesDomain.ErrWorkspaceNotFound.Error()
}
