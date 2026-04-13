package domain

import "errors"

var (
	ErrAssignmentNotFound             = errors.New("assignment not found")
	ErrAssignmentAlreadyExists        = errors.New("assignment already exists for user, workspace and role")
	ErrAssignmentInvalidInput         = errors.New("invalid input")
	ErrAssignmentRoleRequired         = errors.New("assignment role is required")
	ErrAssignmentRoleInvalid          = errors.New("assignment role is invalid")
	ErrAssignmentWeeklyHoursInvalid   = errors.New("assignment weekly hours must be greater than 0")
	ErrAssignmentUserIDRequired       = errors.New("assignment user id is required")
	ErrAssignmentWorkspaceIDRequired = errors.New("assignment workspace id is required")
	ErrAssignmentUserNotFound         = errors.New("assignment user not found")
	ErrAssignmentWorkspaceNotFound    = errors.New("assignment workspace not found")
	ErrAssignmentWorkspaceClosed      = errors.New("assignment workspace is closed")

	ErrAssignmentAssistantHoursLimitExceeded = errors.New("assistant weekly hours cannot exceed 22")
	ErrAssignmentMonitorCountLimitExceeded   = errors.New("monitor assignments cannot exceed 3")
	ErrAssignmentMonitorHoursLimitExceeded   = errors.New("monitor weekly hours cannot exceed 12")
	ErrAssignmentMonitorFortyPercentExceeded = errors.New("monitor weekly hours cannot exceed 40 percent of assistant weekly hours")
)
