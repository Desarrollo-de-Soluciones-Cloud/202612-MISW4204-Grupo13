package domain

import "errors"

var (
	ErrTaskNotFound              = errors.New("task not found")
	ErrInvalidInput              = errors.New("invalid input")
	ErrTaskAssignmentIDRequired  = errors.New("task assignment id is required")
	ErrTaskAssignmentNotFound    = errors.New("task assignment not found")
	ErrTaskWorkspaceNotFound     = errors.New("task workspace not found")
	ErrTaskWorkspaceClosed       = errors.New("task workspace is closed")
	ErrTaskWeekInactive          = errors.New("task week is not active")
	ErrTaskTitleRequired         = errors.New("task title is required")
	ErrTaskDescriptionRequired   = errors.New("task description is required")
	ErrTaskStatusRequired        = errors.New("task status is required")
	ErrTaskStatusInvalid         = errors.New("task status is invalid")
	ErrTaskSpentHoursRequired    = errors.New("task spent hours is required")
	ErrTaskSpentHoursInvalid     = errors.New("task spent hours must be greater than or equal to 1")
	ErrTaskWeekStartDateRequired = errors.New("task week start date is required")
	ErrTaskWeekStartDateInvalid  = errors.New("task week start date is invalid")
	ErrTaskLateUpdateForbidden   = errors.New("late task cannot be updated")
	ErrTaskDeleteForbidden       = errors.New("task can only be deleted during its active week")
	ErrTaskForbidden             = errors.New("forbidden")
)
