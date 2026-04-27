package domain

import "errors"

var (
	ErrReportInvalidInput = errors.New("invalid input")
	ErrReportNotFound = errors.New("report not found")
	ErrReportWorkspaceIDRequired = errors.New("workspace_id is required")
	ErrReportWeekIDRequired = errors.New("week_id is required")
	ErrReportAssignmentIDRequired = errors.New("assignment_id is required")
	ErrReportUserIDRequired = errors.New("user_id is required")
	ErrReportFilePathRequired = errors.New("file_path is required")
	ErrReportWorkspaceFilterRequired = errors.New("workspace_id filter is required")
	ErrReportWorkspaceNotFound = errors.New("workspace not found")
	ErrReportWeekNotFound = errors.New("week not found")
	ErrReportUserNotFound = errors.New("user not found")
	ErrReportWorkspaceAccessDenied = errors.New("workspace does not belong to the current professor")
	ErrReportNoAssignmentsFound = errors.New("workspace has no monitor or assistant assignments")
	ErrReportNoTasksFoundForWeek = errors.New("no tasks found for the selected workspace and week")
	ErrReportAIGenerationFailed = errors.New("ai report generation failed")
	ErrReportPDFGenerationFailed = errors.New("pdf report generation failed")
	ErrReportFileNotFound = errors.New("report file not found")
)