package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	assignmentsDomain "backend/internal/assignments/domain"
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/tasks/application"
	"backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"

	"github.com/gin-gonic/gin"
)

type TaskAssignmentReader interface {
	FindByID(id uint) (*assignmentsDomain.Assignment, error)
}

type TaskWorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type TaskFileStorage interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, contentType string) error
	Download(ctx context.Context, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, objectName string) error
}

const headerContentType = "Content-Type"

type TaskHandlerDependencies struct {
	AssignmentReader  TaskAssignmentReader
	WorkspaceReader   TaskWorkspaceReader
	FileStorage       TaskFileStorage
	AttachmentsPrefix string
}

type TaskHandler struct {
	createTask      *application.CreateTask
	listTasks       *application.ListTasks
	getTaskByID     *application.GetTaskByID
	updateTask      *application.UpdateTask
	setTaskAttachments *application.SetTaskAttachments
	deleteTask      *application.DeleteTask
	assignmentReader TaskAssignmentReader
	workspaceReader  TaskWorkspaceReader
	fileStorage      TaskFileStorage
	attachmentsPrefix string
}

func NewTaskHandler(
	createTask *application.CreateTask,
	listTasks *application.ListTasks,
	getTaskByID *application.GetTaskByID,
	updateTask *application.UpdateTask,
	setTaskAttachments *application.SetTaskAttachments,
	deleteTask *application.DeleteTask,
	deps TaskHandlerDependencies,
) *TaskHandler {
	return &TaskHandler{
		createTask:      createTask,
		listTasks:       listTasks,
		getTaskByID:     getTaskByID,
		updateTask:      updateTask,
		setTaskAttachments: setTaskAttachments,
		deleteTask:      deleteTask,
		assignmentReader: deps.AssignmentReader,
		workspaceReader:  deps.WorkspaceReader,
		fileStorage:      deps.FileStorage,
		attachmentsPrefix: strings.Trim(deps.AttachmentsPrefix, "/"),
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	currentUser, req, existingAttachments, files, weekStartDate, ok := h.prepareTaskWriteRequest(c)
	if !ok {
		return
	}

	if !h.authorizeTaskCreation(c, currentUser, req.AssignmentID) {
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
		Attachments:   existingAttachments,
	})
	if err != nil {
		h.respondCreateTaskError(c, err)
		return
	}

	output, ok = h.appendUploadedAttachments(c, output, files)
	if !ok {
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
	_, id, existing, ok := h.authorizeExistingTaskAccess(c)
	if !ok {
		return
	}

	_, req, existingAttachments, files, weekStartDate, ok := h.prepareTaskWriteRequest(c)
	if !ok {
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
		Attachments:   existingAttachments,
	})
	if err != nil {
		h.respondUpdateTaskError(c, err)
		return
	}

	output, ok = h.appendUploadedAttachments(c, output, files)
	if !ok {
		return
	}

	if err := h.deleteAttachmentsFromBucket(diffRemovedAttachments(existing.Attachments, output.Attachments)); err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
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

	if err := h.deleteAttachmentsFromBucket(existing.Attachments); err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) bindTaskRequest(c *gin.Context) (CreateTaskRequest, []domain.TaskAttachment, []*multipart.FileHeader, error) {
	if strings.Contains(c.GetHeader(headerContentType), "multipart/form-data") {
		return h.bindMultipartTaskRequest(c)
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return CreateTaskRequest{}, nil, nil, err
	}

	return req, []domain.TaskAttachment{}, nil, nil
}

func (h *TaskHandler) bindMultipartTaskRequest(c *gin.Context) (CreateTaskRequest, []domain.TaskAttachment, []*multipart.FileHeader, error) {
	var req CreateTaskRequest

	req.AssignmentID = parseUintField(c.PostForm("assignment_id"))
	req.Title = c.PostForm("title")
	req.Description = c.PostForm("description")
	req.Status = c.PostForm("status")
	req.SpentHours = parseIntField(c.PostForm("spent_hours"))
	req.Observations = c.PostForm("observations")
	req.WeekStartDate = c.PostForm("week_start_date")

	if validationErrors := validateTaskMultipartRequest(req); len(validationErrors) > 0 {
		return CreateTaskRequest{}, nil, nil, validationErrors[0]
	}

	form, err := c.MultipartForm()
	if err != nil {
		return CreateTaskRequest{}, nil, nil, err
	}

	files := form.File["attachments"]
	attachments := make([]domain.TaskAttachment, 0)

	if rawExisting := c.PostForm("existing_attachments"); strings.TrimSpace(rawExisting) != "" {
		var existing []domain.TaskAttachment
		if err := json.Unmarshal([]byte(rawExisting), &existing); err != nil {
			return CreateTaskRequest{}, nil, nil, domain.ErrInvalidInput
		}
		attachments = append(existing, attachments...)
	}

	return req, attachments, files, nil
}

func (h *TaskHandler) uploadAttachments(taskID uint, startIndex int, files []*multipart.FileHeader) ([]domain.TaskAttachment, error) {
	attachments := make([]domain.TaskAttachment, 0, len(files))
	for index, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		attachment, uploadErr := h.uploadAttachment(taskID, startIndex+index, fileHeader, file)
		_ = file.Close()
		if uploadErr != nil {
			return nil, uploadErr
		}

		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

func (h *TaskHandler) uploadAttachment(
	taskID uint,
	fileNumber int,
	fileHeader *multipart.FileHeader,
	file io.Reader,
) (domain.TaskAttachment, error) {
	if h.fileStorage == nil {
		return domain.TaskAttachment{}, sharedErrors.ErrInternalServerError
	}

	fileName := sanitizeFileName(fileHeader.Filename)
	objectName := buildAttachmentObjectName(h.attachmentsPrefix, taskID, fileNumber, fileName)
	contentType := fileHeader.Header.Get(headerContentType)
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	if err := h.fileStorage.Upload(context.Background(), objectName, file, contentType); err != nil {
		return domain.TaskAttachment{}, err
	}

	return domain.TaskAttachment{
		ID:          buildAttachmentID(),
		Name:        fileHeader.Filename,
		FilePath:    objectName,
		ContentType: contentType,
		Size:        fileHeader.Size,
	}, nil
}

func (h *TaskHandler) DownloadAttachment(c *gin.Context) {
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

	task, err := h.getTaskByID.Execute(application.GetTaskByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound), errors.Is(err, domain.ErrTaskAssignmentNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	if !h.canAccessTask(currentUser.GlobalRole, currentUser.ID, task.UserID, task.AssignmentID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	attachmentID := strings.TrimSpace(c.Param("attachmentId"))
	attachment, found := findAttachmentByID(task.Attachments, attachmentID)
	if !found {
		sharedHelpers.RespondWithError(c, http.StatusNotFound, sharedErrors.ErrNotFound)
		return
	}
	reader, err := h.fileStorage.Download(context.Background(), attachment.FilePath)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusNotFound, sharedErrors.ErrNotFound)
		return
	}
	defer reader.Close()

	contentType := attachment.ContentType
	if strings.TrimSpace(contentType) == "" {
		contentType = "application/octet-stream"
	}

	c.Header(headerContentType, contentType)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(attachment.Name)))

	if _, err := io.Copy(c.Writer, reader); err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
	}
}

func (h *TaskHandler) prepareTaskWriteRequest(c *gin.Context) (authDomain.AuthenticatedUser, CreateTaskRequest, []domain.TaskAttachment, []*multipart.FileHeader, time.Time, bool) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return authDomain.AuthenticatedUser{}, CreateTaskRequest{}, nil, nil, time.Time{}, false
	}

	req, existingAttachments, files, err := h.bindTaskRequest(c)
	if err != nil {
		h.respondTaskBindingError(c, err)
		return authDomain.AuthenticatedUser{}, CreateTaskRequest{}, nil, nil, time.Time{}, false
	}

	weekStartDate, err := parseWeekStartDate(req.WeekStartDate)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return authDomain.AuthenticatedUser{}, CreateTaskRequest{}, nil, nil, time.Time{}, false
	}

	return currentUser, req, existingAttachments, files, weekStartDate, true
}

func (h *TaskHandler) authorizeTaskCreation(c *gin.Context, currentUser authDomain.AuthenticatedUser, assignmentID uint) bool {
	if currentUser.GlobalRole == usersDomain.RoleAdmin {
		return true
	}

	assignment, err := h.assignmentReader.FindByID(assignmentID)
	if err == nil && !h.canAccessTask(currentUser.GlobalRole, currentUser.ID, assignment.UserID, assignmentID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return false
	}

	return true
}

func (h *TaskHandler) authorizeExistingTaskAccess(c *gin.Context) (authDomain.AuthenticatedUser, uint, *application.TaskOutput, bool) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return authDomain.AuthenticatedUser{}, 0, nil, false
	}

	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return authDomain.AuthenticatedUser{}, 0, nil, false
	}

	existing, err := h.getTaskByID.Execute(application.GetTaskByIDInput{ID: id})
	if err != nil {
		h.respondTaskLookupError(c, err)
		return authDomain.AuthenticatedUser{}, 0, nil, false
	}

	if !h.canAccessTask(currentUser.GlobalRole, currentUser.ID, existing.UserID, existing.AssignmentID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return authDomain.AuthenticatedUser{}, 0, nil, false
	}

	return currentUser, id, existing, true
}

func (h *TaskHandler) appendUploadedAttachments(c *gin.Context, output *application.TaskOutput, files []*multipart.FileHeader) (*application.TaskOutput, bool) {
	if len(files) == 0 {
		return output, true
	}

	attachments, uploadErr := h.uploadAttachments(output.ID, len(output.Attachments)+1, files)
	if uploadErr != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return nil, false
	}

	updatedOutput, err := h.setTaskAttachments.Execute(application.SetTaskAttachmentsInput{
		ID:          output.ID,
		Attachments: append(output.Attachments, attachments...),
	})
	if err != nil {
		_ = h.deleteAttachmentsFromBucket(attachments)
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return nil, false
	}

	return updatedOutput, true
}

func (h *TaskHandler) respondTaskBindingError(c *gin.Context, err error) {
	if isTaskValidationError(err) {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
}

func (h *TaskHandler) respondTaskLookupError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound), errors.Is(err, domain.ErrTaskAssignmentNotFound):
		sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
	default:
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
	}
}

func (h *TaskHandler) respondCreateTaskError(c *gin.Context, err error) {
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
}

func (h *TaskHandler) respondUpdateTaskError(c *gin.Context, err error) {
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

func parseUintField(value string) uint {
	parsed, _ := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	return uint(parsed)
}

func parseIntField(value string) int {
	parsed, _ := strconv.Atoi(strings.TrimSpace(value))
	return parsed
}

func buildAttachmentID() string {
	return fmt.Sprintf("att_%d", time.Now().UnixNano())
}

func validateTaskMultipartRequest(req CreateTaskRequest) []error {
	result := make([]error, 0)
	if req.AssignmentID == 0 {
		result = append(result, domain.ErrTaskAssignmentIDRequired)
	}
	if strings.TrimSpace(req.Title) == "" {
		result = append(result, domain.ErrTaskTitleRequired)
	}
	if strings.TrimSpace(req.Description) == "" {
		result = append(result, domain.ErrTaskDescriptionRequired)
	}
	if strings.TrimSpace(req.Status) == "" {
		result = append(result, domain.ErrTaskStatusRequired)
	}
	if req.SpentHours == 0 {
		result = append(result, domain.ErrTaskSpentHoursRequired)
	}
	if strings.TrimSpace(req.WeekStartDate) == "" {
		result = append(result, domain.ErrTaskWeekStartDateRequired)
	}
	return result
}

func sanitizeFileName(fileName string) string {
	fileName = filepath.Base(strings.ReplaceAll(fileName, "\\", "/"))
	fileName = strings.ReplaceAll(fileName, " ", "_")
	if fileName == "." || fileName == "" {
		return "attachment"
	}
	return fileName
}

func buildAttachmentObjectName(prefix string, taskID uint, fileNumber int, fileName string) string {
	safePrefix := strings.Trim(prefix, "/")
	folder := fmt.Sprintf("task_%d", taskID)
	objectName := fmt.Sprintf("file_%d_%s", fileNumber, fileName)
	if safePrefix == "" {
		return fmt.Sprintf("%s/%s", folder, objectName)
	}
	return fmt.Sprintf("%s/%s/%s", safePrefix, folder, objectName)
}

func findAttachmentByID(attachments []domain.TaskAttachment, attachmentID string) (domain.TaskAttachment, bool) {
	for _, attachment := range attachments {
		if attachment.ID == attachmentID {
			return attachment, true
		}
	}

	return domain.TaskAttachment{}, false
}

func diffRemovedAttachments(before []domain.TaskAttachment, after []domain.TaskAttachment) []domain.TaskAttachment {
	if len(before) == 0 {
		return []domain.TaskAttachment{}
	}

	current := make(map[string]struct{}, len(after))
	for _, attachment := range after {
		current[attachment.FilePath] = struct{}{}
	}

	removed := make([]domain.TaskAttachment, 0)
	for _, attachment := range before {
		if _, ok := current[attachment.FilePath]; ok {
			continue
		}
		removed = append(removed, attachment)
	}

	return removed
}

func (h *TaskHandler) deleteAttachmentsFromBucket(attachments []domain.TaskAttachment) error {
	if h.fileStorage == nil || len(attachments) == 0 {
		return nil
	}

	for _, attachment := range attachments {
		if strings.TrimSpace(attachment.FilePath) == "" {
			continue
		}
		if err := h.fileStorage.Delete(context.Background(), attachment.FilePath); err != nil {
			return err
		}
	}

	return nil
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
		Attachments:   toTaskAttachmentResponses(task.Attachments),
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
