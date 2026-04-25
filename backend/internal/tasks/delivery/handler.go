package delivery

import (
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	assignmentsDomain "backend/internal/assignments/domain"
	"backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskAssignmentReader interface {
	FindByID(id uint) (*assignmentsDomain.Assignment, error)
}

type TaskWorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type TaskHandler struct {
	createTask      *application.CreateTask
	listTasks       *application.ListTasks
	getTaskByID     *application.GetTaskByID
	updateTask      *application.UpdateTask
	deleteTask      *application.DeleteTask
	assignmentReader TaskAssignmentReader
	workspaceReader  TaskWorkspaceReader
}

func NewTaskHandler(
	createTask *application.CreateTask,
	listTasks *application.ListTasks,
	getTaskByID *application.GetTaskByID,
	updateTask *application.UpdateTask,
	deleteTask *application.DeleteTask,
	assignmentReader TaskAssignmentReader,
	workspaceReader TaskWorkspaceReader,
) *TaskHandler {
	return &TaskHandler{
		createTask:      createTask,
		listTasks:       listTasks,
		getTaskByID:     getTaskByID,
		updateTask:      updateTask,
		deleteTask:      deleteTask,
		assignmentReader: assignmentReader,
		workspaceReader:  workspaceReader,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	if currentUser.GlobalRole != usersDomain.RoleAdmin {
		assignment, err := h.assignmentReader.FindByID(req.AssignmentID)
		if err == nil && !h.canAccessTask(currentUser.GlobalRole, currentUser.ID, assignment.UserID, req.AssignmentID) {
			sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
			return
		}
	}

	weekStartDate, err := parseWeekStartDate(req.WeekStartDate)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.createTask.Execute(application.CreateTaskInput{
		AssignmentID:  req.AssignmentID,
		Title:         req.Title,
		Description:   req.Description,
		Status:        domain.TaskStatus(req.Status),
		SpentHours:    req.SpentHours,
		Observations:  req.Observations,
		WeekStartDate: weekStartDate,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskAssignmentNotFound),
			errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskWorkspaceNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrTaskWorkspaceClosed),
			errors.Is(err, domain.ErrTaskWeekInactive):
			sharedHelpers.RespondWithError(c, http.StatusConflict, err)
		case isTaskValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, toTaskResponse(output))
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	output, err := h.listTasks.Execute()
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	tasks := make([]TaskResponse, 0, len(output.Tasks))
	for i := range output.Tasks {
		t := &output.Tasks[i]
		if h.canAccessTask(currentUser.GlobalRole, currentUser.ID, t.UserID, t.AssignmentID) {
			tasks = append(tasks, toTaskResponse(t))
		}
	}

	c.JSON(http.StatusOK, ListTasksResponse{Tasks: tasks})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
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

	output, err := h.getTaskByID.Execute(application.GetTaskByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	if !h.canAccessTask(currentUser.GlobalRole, currentUser.ID, output.UserID, output.AssignmentID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(output))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
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

	existing, err := h.getTaskByID.Execute(application.GetTaskByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	if !h.canAccessTask(currentUser.GlobalRole, currentUser.ID, existing.UserID, existing.AssignmentID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	weekStartDate, err := parseWeekStartDate(req.WeekStartDate)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.updateTask.Execute(application.UpdateTaskInput{
		ID:            id,
		AssignmentID:  req.AssignmentID,
		Title:         req.Title,
		Description:   req.Description,
		Status:        domain.TaskStatus(req.Status),
		SpentHours:    req.SpentHours,
		Observations:  req.Observations,
		WeekStartDate: weekStartDate,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case isTaskValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		case errors.Is(err, domain.ErrTaskLateUpdateForbidden):
			sharedHelpers.RespondWithError(c, http.StatusForbidden, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(output))
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
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

	existing, err := h.getTaskByID.Execute(application.GetTaskByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	if !h.canAccessTask(currentUser.GlobalRole, currentUser.ID, existing.UserID, existing.AssignmentID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	err = h.deleteTask.Execute(application.DeleteTaskInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrTaskDeleteForbidden):
			sharedHelpers.RespondWithError(c, http.StatusForbidden, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func isTaskValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrTaskAssignmentIDRequired) ||
		errors.Is(err, domain.ErrTaskTitleRequired) ||
		errors.Is(err, domain.ErrTaskDescriptionRequired) ||
		errors.Is(err, domain.ErrTaskStatusRequired) ||
		errors.Is(err, domain.ErrTaskStatusInvalid) ||
		errors.Is(err, domain.ErrTaskSpentHoursRequired) ||
		errors.Is(err, domain.ErrTaskSpentHoursInvalid) ||
		errors.Is(err, domain.ErrTaskWeekStartDateRequired) ||
		errors.Is(err, domain.ErrTaskWeekStartDateInvalid)
}

func toTaskResponse(task *application.TaskOutput) TaskResponse {
	return TaskResponse{
		ID:            task.ID,
		UserID:        task.UserID,
		AssignmentID:  task.AssignmentID,
		WeekID:        task.WeekID,
		Title:         task.Title,
		Description:   task.Description,
		Status:        string(task.Status),
		SpentHours:    task.SpentHours,
		Observations:  task.Observations,
		WeekStartDate: task.WeekStartDate.Format(dateLayout),
		Late:          task.Late,
	}
}

func (h *TaskHandler) canAccessTask(role usersDomain.UserRole, currentUserID, taskUserID, assignmentID uint) bool {
	switch role {
	case usersDomain.RoleAdmin:
		return true
	case usersDomain.RoleMonitor, usersDomain.RoleAssistant:
		return currentUserID == taskUserID
	case usersDomain.RoleProfessor:
		assignment, err := h.assignmentReader.FindByID(assignmentID)
		if err != nil {
			return false
		}
		workspace, err := h.workspaceReader.FindByID(assignment.WorkspaceID)
		if err != nil {
			return false
		}
		return workspace.UserID == currentUserID
	default:
		return false
	}
}
