package delivery

import (
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	createTask  *application.CreateTask
	listTasks   *application.ListTasks
	getTaskByID *application.GetTaskByID
	updateTask  *application.UpdateTask
	deleteTask  *application.DeleteTask
}

func NewTaskHandler(
	createTask *application.CreateTask,
	listTasks *application.ListTasks,
	getTaskByID *application.GetTaskByID,
	updateTask *application.UpdateTask,
	deleteTask *application.DeleteTask,
) *TaskHandler {
	return &TaskHandler{
		createTask:  createTask,
		listTasks:   listTasks,
		getTaskByID: getTaskByID,
		updateTask:  updateTask,
		deleteTask:  deleteTask,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	output, err := h.createTask.Execute(application.CreateTaskInput{
		AssignmentID: req.AssignmentID,
		WeekID:       req.WeekID,
		Title:        req.Title,
		Description:  req.Description,
		Status:       domain.TaskStatus(req.Status),
		SpentHours:   req.SpentHours,
		Observations: req.Observations,
		Attachments:  toAttachmentInputs(req.Attachments),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskAssignmentNotFound),
			errors.Is(err, domain.ErrTaskWorkspaceNotFound),
			errors.Is(err, domain.ErrTaskWeekNotFound),
			errors.Is(err, domain.ErrTaskNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
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
	output, err := h.listTasks.Execute()
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	tasks := make([]TaskResponse, len(output.Tasks))
	for i := range output.Tasks {
		tasks[i] = toTaskResponse(&output.Tasks[i])
	}

	c.JSON(http.StatusOK, ListTasksResponse{Tasks: tasks})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
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

	c.JSON(http.StatusOK, toTaskResponse(output))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	output, err := h.updateTask.Execute(application.UpdateTaskInput{
		ID:           id,
		AssignmentID: req.AssignmentID,
		WeekID:       req.WeekID,
		Title:        req.Title,
		Description:  req.Description,
		Status:       domain.TaskStatus(req.Status),
		SpentHours:   req.SpentHours,
		Observations: req.Observations,
		Attachments:  toAttachmentInputs(req.Attachments),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskWorkspaceNotFound),
			errors.Is(err, domain.ErrTaskWeekNotFound),
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
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	err = h.deleteTask.Execute(application.DeleteTaskInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound),
			errors.Is(err, domain.ErrTaskWeekNotFound),
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
		errors.Is(err, domain.ErrTaskAssignmentChangeForbidden) ||
		errors.Is(err, domain.ErrTaskWeekIDRequired) ||
		errors.Is(err, domain.ErrTaskWeekChangeForbidden) ||
		errors.Is(err, domain.ErrTaskTitleRequired) ||
		errors.Is(err, domain.ErrTaskDescriptionRequired) ||
		errors.Is(err, domain.ErrTaskStatusRequired) ||
		errors.Is(err, domain.ErrTaskStatusInvalid) ||
		errors.Is(err, domain.ErrTaskSpentHoursRequired) ||
		errors.Is(err, domain.ErrTaskSpentHoursInvalid) ||
		errors.Is(err, domain.ErrTaskWeekPeriodMismatch) ||
		errors.Is(err, domain.ErrTaskAttachmentPathRequired) ||
		errors.Is(err, domain.ErrTaskWeekStartDateRequired) ||
		errors.Is(err, domain.ErrTaskWeekStartDateInvalid)
}

func toTaskResponse(task *application.TaskOutput) TaskResponse {
	attachments := make([]TaskAttachmentResponse, len(task.Attachments))
	for i, attachment := range task.Attachments {
		attachments[i] = TaskAttachmentResponse{
			ID:   attachment.ID,
			Path: attachment.Path,
		}
	}

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
		Attachments:   attachments,
	}
}

func toAttachmentInputs(attachments []TaskAttachmentRequest) []application.TaskAttachmentInput {
	result := make([]application.TaskAttachmentInput, len(attachments))
	for i, attachment := range attachments {
		result[i] = application.TaskAttachmentInput{Path: attachment.Path}
	}
	return result
}
