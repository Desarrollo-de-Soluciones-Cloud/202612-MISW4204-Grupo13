package delivery

import (
	"backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
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
	userReader              AssignmentUserReader
}

type AssignmentWorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type AssignmentUserReader interface {
	FindByID(id uint) (*usersDomain.User, error)
}

func NewAssignmentHandler(
	createAssignment *application.CreateAssignment,
	getAssignmentByID *application.GetAssignmentByID,
	listAssignmentsByUserID *application.ListAssignmentsByUserID,
	updateAssignment *application.UpdateAssignment,
	workspaceReader AssignmentWorkspaceReader,
	userReader AssignmentUserReader,
) *AssignmentHandler {
	return &AssignmentHandler{
		createAssignment:        createAssignment,
		getAssignmentByID:       getAssignmentByID,
		listAssignmentsByUserID: listAssignmentsByUserID,
		updateAssignment:        updateAssignment,
		workspaceReader:         workspaceReader,
		userReader:              userReader,
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

	// Additional validation: Professors can only assign to monitor or assistant users
	var userGlobalRole usersDomain.UserRole
	if currentUser.GlobalRole == usersDomain.RoleProfessor {
		// We need to fetch the user to validate their global role
		role, err := h.validateUserRoleForAssignment(req.UserID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrAssignmentUserNotFound):
				sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
			case errors.Is(err, domain.ErrAssignmentUserInvalidRole):
				sharedHelpers.RespondWithError(c, http.StatusForbidden, err)
			default:
				sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
			}
			return
		}
		userGlobalRole = role
	} else {
		// For admin, we still need to get the user's global role
		role, err := h.validateUserRoleForAssignment(req.UserID)
		if err != nil {
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
			return
		}
		userGlobalRole = role
	}

	// Professors can only assign to monitor or assistant users
	if currentUser.GlobalRole == usersDomain.RoleProfessor {
		if userGlobalRole != usersDomain.RoleMonitor && userGlobalRole != usersDomain.RoleAssistant {
			sharedHelpers.RespondWithError(c, http.StatusForbidden, domain.ErrAssignmentUserInvalidRole)
			return
		}
	}

	// Validate assignment role is compatible with user's global role
	// If user is monitor, assignment role must be monitor
	if userGlobalRole == usersDomain.RoleMonitor {
		if req.Role != string(domain.RoleMonitor) {
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, domain.ErrAssignmentRoleNotAllowedForUser)
			return
		}
	}
	// If user is assistant, assignment role can be monitor or assistant (no additional check needed)

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
			errors.Is(err, domain.ErrAssignmentPeriodClosed),
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

	var req UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	// Only Admin and Professor can update assignments
	if currentUser.GlobalRole != usersDomain.RoleAdmin && currentUser.GlobalRole != usersDomain.RoleProfessor {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	// If Professor, get assignment and workspace to validate permissions
	if currentUser.GlobalRole == usersDomain.RoleProfessor {
		// Get current assignment
		currentAssignment, err := h.getAssignmentByID.Execute(application.GetAssignmentByIDInput{ID: id})
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrAssignmentNotFound):
				sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
			default:
				sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
			}
			return
		}

		// Get workspace to verify professor ownership
		workspace, err := h.workspaceReader.FindByID(currentAssignment.WorkspaceID)
		if err != nil {
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
			return
		}

		// Check if professor owns the workspace
		if workspace.UserID != currentUser.ID {
			sharedHelpers.RespondWithError(c, http.StatusForbidden, domain.ErrAssignmentProfessorCannotUpdate)
			return
		}

		// Check if professor is trying to change weekly hours
		if currentAssignment.WeeklyHours != req.WeeklyHours {
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, domain.ErrAssignmentProfessorCannotChangeWeeklyHours)
			return
		}
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
		case errors.Is(err, domain.ErrAssignmentWorkspaceClosed),
			errors.Is(err, domain.ErrAssignmentPeriodClosed),
			errors.Is(err, domain.ErrAssignmentAlreadyExists):
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

func (h *AssignmentHandler) validateUserRoleForAssignment(userID uint) (usersDomain.UserRole, error) {
	if h.userReader == nil {
		// If no user reader, we can't validate - this shouldn't happen
		return "", errors.New("user reader not initialized")
	}

	user, err := h.userReader.FindByID(userID)
	if err != nil {
		return "", domain.ErrAssignmentUserNotFound
	}

	return user.GlobalRole, nil
}

func isWorkspaceNotFoundError(err error) bool {
	return errors.Is(err, workspacesDomain.ErrWorkspaceNotFound) || err.Error() == workspacesDomain.ErrWorkspaceNotFound.Error()
}
