package delivery

import (
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/assignments/application"
	"backend/internal/assignments/domain"
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
}

func NewAssignmentHandler(
	createAssignment *application.CreateAssignment,
	getAssignmentByID *application.GetAssignmentByID,
	listAssignmentsByUserID *application.ListAssignmentsByUserID,
	updateAssignment *application.UpdateAssignment,
) *AssignmentHandler {
	return &AssignmentHandler{
		createAssignment:        createAssignment,
		getAssignmentByID:       getAssignmentByID,
		listAssignmentsByUserID: listAssignmentsByUserID,
		updateAssignment:        updateAssignment,
	}
}

func (h *AssignmentHandler) CreateAssignment(c *gin.Context) {
	var req CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
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
	rawUserID := c.Query("user_id")
	parsedID, err := strconv.ParseUint(rawUserID, 10, 32)
	if err != nil || parsedID == 0 {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, domain.ErrAssignmentUserIDRequired)
		return
	}

	output, err := h.listAssignmentsByUserID.Execute(application.ListAssignmentsByUserIDInput{
		UserID: uint(parsedID),
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

	assignments := make([]AssignmentResponse, len(output.Assignments))
	for i, a := range output.Assignments {
		assignments[i] = AssignmentResponse{
			ID:          a.ID,
			UserID:      a.UserID,
			WorkspaceID: a.WorkspaceID,
			Role:        string(a.Role),
			WeeklyHours: a.WeeklyHours,
		}
	}

	c.JSON(http.StatusOK, ListAssignmentsResponse{Assignments: assignments})
}
