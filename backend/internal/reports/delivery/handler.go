package delivery

import (
	"encoding/base64"
	"encoding/json"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	reportsApplication "backend/internal/reports/application"
	reportsDomain "backend/internal/reports/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"

	"github.com/gin-gonic/gin"
)

type ReportWorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type ReportFileStorage interface {
	Download(ctx context.Context, objectName string) (io.ReadCloser, error)
}

type ReportHandler struct {
	generateWeeklyReports *reportsApplication.GenerateWeeklyReports
	queueWeeklyReports    *reportsApplication.QueueWeeklyReports
	processWeeklyReport   *reportsApplication.ProcessWeeklyReportJob
	listReports           *reportsApplication.ListReports
	getReportByID         *reportsApplication.GetReportByID
	workspaceReader       ReportWorkspaceReader
	reportFileStorage     ReportFileStorage
	pubSubPushAuthToken   string
}

func NewReportHandler(
	generateWeeklyReports *reportsApplication.GenerateWeeklyReports,
	queueWeeklyReports *reportsApplication.QueueWeeklyReports,
	processWeeklyReport *reportsApplication.ProcessWeeklyReportJob,
	listReports *reportsApplication.ListReports,
	getReportByID *reportsApplication.GetReportByID,
	workspaceReader ReportWorkspaceReader,
	reportFileStorage ReportFileStorage,
	pubSubPushAuthToken string,
) *ReportHandler {
	return &ReportHandler{
		generateWeeklyReports: generateWeeklyReports,
		queueWeeklyReports:    queueWeeklyReports,
		processWeeklyReport:   processWeeklyReport,
		listReports:           listReports,
		getReportByID:         getReportByID,
		workspaceReader:       workspaceReader,
		reportFileStorage:     reportFileStorage,
		pubSubPushAuthToken:   pubSubPushAuthToken,
	}
}

func (h *ReportHandler) GenerateWeeklyReports(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	var req GenerateWeeklyReportsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	workspace, err := h.workspaceReader.FindByID(req.WorkspaceID)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusNotFound, reportsDomain.ErrReportWorkspaceNotFound)
		return
	}

	if !canAccessWorkspace(currentUser, workspace.UserID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, reportsDomain.ErrReportWorkspaceAccessDenied)
		return
	}

	if h.queueWeeklyReports != nil {
		h.queueReports(c, req)
		return
	}

	output, err := h.generateWeeklyReports.Execute(reportsApplication.GenerateWeeklyReportsInput{
		WorkspaceID: req.WorkspaceID,
		WeekID:      req.WeekID,
	})
	if err != nil {
		handleGenerateReportsError(c, err)
		return
	}

	responses := make([]ReportResponse, 0, len(output.Reports))
	for _, report := range output.Reports {
		responses = append(responses, toReportResponse(report))
	}

	c.JSON(http.StatusCreated, GenerateWeeklyReportsResponse{
		Reports:        responses,
		GeneratedCount: len(responses),
	})
}

func (h *ReportHandler) ProcessWeeklyReportJob(c *gin.Context) {
	if !h.isAuthorizedPubSubPushRequest(c) {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	if h.processWeeklyReport == nil {
		sharedHelpers.RespondWithError(c, http.StatusNotImplemented, sharedErrors.ErrInternalServerError)
		return
	}

	job, err := parsePubSubWeeklyReportJob(c)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, reportsDomain.ErrReportInvalidInput)
		return
	}

	if _, err := h.processWeeklyReport.Execute(job); err != nil {
		handleGenerateReportsError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReportHandler) ListReports(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	workspaceID, err := parseRequiredWorkspaceID(c.Query("workspace_id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	weekID, err := parseOptionalResourceID(c.Query("week_id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	userID, err := parseOptionalResourceID(c.Query("user_id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	workspace, err := h.workspaceReader.FindByID(workspaceID)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusNotFound, reportsDomain.ErrReportWorkspaceNotFound)
		return
	}

	if !canAccessWorkspace(currentUser, workspace.UserID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, reportsDomain.ErrReportWorkspaceAccessDenied)
		return
	}

	output, err := h.listReports.Execute(reportsApplication.ListReportsInput{
		WorkspaceID: workspaceID,
		WeekID:      weekID,
		UserID:      userID,
	})
	if err != nil {
		handleListReportsError(c, err)
		return
	}

	responses := make([]ReportResponse, 0, len(output.Reports))
	for _, report := range output.Reports {
		responses = append(responses, toReportResponse(report))
	}

	c.JSON(http.StatusOK, ListReportsResponse{Reports: responses})
}

func (h *ReportHandler) DownloadReport(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	reportID, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	report, err := h.getReportByID.Execute(reportID)
	if err != nil {
		if errors.Is(err, reportsDomain.ErrReportNotFound) {
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
			return
		}

		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	if !h.canAccessReport(currentUser, report.WorkspaceID) {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, reportsDomain.ErrReportWorkspaceAccessDenied)
		return
	}

	if h.reportFileStorage == nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	reader, err := h.reportFileStorage.Download(context.Background(), report.FilePath)
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusNotFound, reportsDomain.ErrReportFileNotFound)
		return
	}
	defer reader.Close()

	fileName := filepath.Base(strings.ReplaceAll(report.FilePath, "\\", "/"))

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))

	if _, err := io.Copy(c.Writer, reader); err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}
}

func handleGenerateReportsError(c *gin.Context, err error) {
	switch {
	case isReportValidationError(err),
		errors.Is(err, reportsDomain.ErrReportNoAssignmentsFound),
		errors.Is(err, reportsDomain.ErrReportNoTasksFoundForWeek):
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)

	case errors.Is(err, reportsDomain.ErrReportWorkspaceNotFound),
		errors.Is(err, reportsDomain.ErrReportWeekNotFound):
		sharedHelpers.RespondWithError(c, http.StatusNotFound, err)

	case errors.Is(err, reportsDomain.ErrReportWorkspaceAccessDenied):
		sharedHelpers.RespondWithError(c, http.StatusForbidden, err)

	case errors.Is(err, reportsDomain.ErrReportAIGenerationFailed),
		errors.Is(err, reportsDomain.ErrReportPDFGenerationFailed):
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, err)

	default:
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
	}
}

func handleListReportsError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, reportsDomain.ErrReportWorkspaceFilterRequired):
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)

	case errors.Is(err, reportsDomain.ErrReportWorkspaceNotFound),
		errors.Is(err, reportsDomain.ErrReportWeekNotFound),
		errors.Is(err, reportsDomain.ErrReportUserNotFound):
		sharedHelpers.RespondWithError(c, http.StatusNotFound, err)

	case errors.Is(err, reportsDomain.ErrReportWorkspaceAccessDenied):
		sharedHelpers.RespondWithError(c, http.StatusForbidden, err)

	default:
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
	}
}

func (h *ReportHandler) canAccessReport(currentUser authDomain.AuthenticatedUser, workspaceID uint) bool {
	if currentUser.GlobalRole == usersDomain.RoleAdmin {
		return true
	}

	if currentUser.GlobalRole != usersDomain.RoleProfessor {
		return false
	}

	workspace, err := h.workspaceReader.FindByID(workspaceID)
	if err != nil {
		return false
	}

	return workspace.UserID == currentUser.ID
}

func canAccessWorkspace(currentUser authDomain.AuthenticatedUser, workspaceOwnerID uint) bool {
	if currentUser.GlobalRole == usersDomain.RoleAdmin {
		return true
	}

	if currentUser.GlobalRole == usersDomain.RoleProfessor {
		return currentUser.ID == workspaceOwnerID
	}

	return false
}

func toReportResponse(report reportsApplication.ReportOutput) ReportResponse {
	return ReportResponse{
		ID:           report.ID,
		WorkspaceID:  report.WorkspaceID,
		WeekID:       report.WeekID,
		AssignmentID: report.AssignmentID,
		UserID:       report.UserID,
		FilePath:     report.FilePath,
	}
}

func (h *ReportHandler) queueReports(c *gin.Context, req GenerateWeeklyReportsRequest) {
	output, err := h.queueWeeklyReports.Execute(reportsApplication.QueueWeeklyReportsInput{
		WorkspaceID: req.WorkspaceID,
		WeekID:      req.WeekID,
	})
	if err != nil {
		handleGenerateReportsError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, GenerateWeeklyReportsResponse{
		Reports:        []ReportResponse{},
		GeneratedCount: output.QueuedCount,
	})
}

func (h *ReportHandler) isAuthorizedPubSubPushRequest(c *gin.Context) bool {
	if strings.TrimSpace(h.pubSubPushAuthToken) == "" {
		return false
	}

	return c.GetHeader("X-PubSub-Token") == h.pubSubPushAuthToken ||
		c.Query("token") == h.pubSubPushAuthToken
}

func parsePubSubWeeklyReportJob(c *gin.Context) (reportsApplication.WeeklyReportJobMessage, error) {
	var req PubSubPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return reportsApplication.WeeklyReportJobMessage{}, err
	}

	decoded, err := base64.StdEncoding.DecodeString(req.Message.Data)
	if err != nil {
		return reportsApplication.WeeklyReportJobMessage{}, err
	}

	var job reportsApplication.WeeklyReportJobMessage
	if err := json.Unmarshal(decoded, &job); err != nil {
		return reportsApplication.WeeklyReportJobMessage{}, err
	}

	return job, nil
}
