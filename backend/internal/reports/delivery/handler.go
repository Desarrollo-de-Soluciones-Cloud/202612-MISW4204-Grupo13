package delivery

import (
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	reportsApplication "backend/internal/reports/application"
	reportsDomain "backend/internal/reports/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"

	"github.com/gin-gonic/gin"
)

type ReportWorkspaceReader interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type ReportHandler struct {
	generateWeeklyReports *reportsApplication.GenerateWeeklyReports
	listReports           *reportsApplication.ListReports
	getReportByID         *reportsApplication.GetReportByID
	workspaceReader       ReportWorkspaceReader
}

func NewReportHandler(
	generateWeeklyReports *reportsApplication.GenerateWeeklyReports,
	listReports *reportsApplication.ListReports,
	getReportByID *reportsApplication.GetReportByID,
	workspaceReader ReportWorkspaceReader,
) *ReportHandler {
	return &ReportHandler{
		generateWeeklyReports: generateWeeklyReports,
		listReports:           listReports,
		getReportByID:         getReportByID,
		workspaceReader:       workspaceReader,
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
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	output, err := h.generateWeeklyReports.Execute(reportsApplication.GenerateWeeklyReportsInput{
		WorkspaceID: req.WorkspaceID,
		WeekID:      req.WeekID,
	})
	if err != nil {
		switch {
		case errors.Is(err, reportsDomain.ErrReportWorkspaceNotFound), errors.Is(err, reportsDomain.ErrReportWeekNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, reportsDomain.ErrReportWorkspaceIDRequired), errors.Is(err, reportsDomain.ErrReportWeekIDRequired):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
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

func (h *ReportHandler) ListReports(c *gin.Context) {
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	workspaceID, err := parseOptionalUint(c.Query("workspace_id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	weekID, err := parseOptionalUint(c.Query("week_id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.listReports.Execute(reportsApplication.ListReportsInput{
		WorkspaceID: workspaceID,
		WeekID:      weekID,
	})
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		return
	}

	responses := make([]ReportResponse, 0, len(output.Reports))
	for _, report := range output.Reports {
		if !h.canAccessReport(currentUser, report.WorkspaceID) {
			continue
		}
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
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	if _, err := os.Stat(report.FilePath); err != nil {
		sharedHelpers.RespondWithError(c, http.StatusNotFound, reportsDomain.ErrReportFileNotFound)
		return
	}

	c.FileAttachment(report.FilePath, filepath.Base(report.FilePath))
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
		Type:         report.Type,
		Summary:      report.Summary,
		FilePath:     report.FilePath,
	}
}
