package domain

import "errors"

var (
	ErrTaskNotFound                  = errors.New("task not found")
	ErrInvalidInput                  = errors.New("invalid input")
	ErrTaskAssignmentIDRequired      = errors.New("task assignment id is required")
	ErrTaskAssignmentNotFound        = errors.New("task assignment not found")
	ErrTaskAssignmentChangeForbidden = errors.New("task assignment cannot be changed once created")
	ErrTaskWorkspaceNotFound         = errors.New("task workspace not found")
	ErrTaskWeekIDRequired            = errors.New("task week id is required")
	ErrTaskWeekNotFound              = errors.New("task week not found")
	ErrTaskWeekChangeForbidden       = errors.New("task week cannot be changed once created")
	ErrTaskWeekPeriodMismatch        = errors.New("task week does not belong to assignment workspace period")
	ErrTaskTitleRequired             = errors.New("task title is required")
	ErrTaskDescriptionRequired       = errors.New("task description is required")
	ErrTaskStatusRequired            = errors.New("task status is required")
	ErrTaskStatusInvalid             = errors.New("task status is invalid")
	ErrTaskSpentHoursRequired        = errors.New("task spent hours is required")
	ErrTaskSpentHoursInvalid         = errors.New("task spent hours must be greater than or equal to 1")
	ErrTaskWeekStartDateRequired     = errors.New("task week start date is required")
	ErrTaskWeekStartDateInvalid      = errors.New("task week start date is invalid")
	ErrTaskAttachmentPathRequired    = errors.New("task attachment path is required")
	ErrTaskLateUpdateForbidden       = errors.New("late task cannot be updated")
	ErrTaskDeleteForbidden           = errors.New("task can only be deleted during its active week")
	ErrTaskForbidden                 = errors.New("forbidden")
)
