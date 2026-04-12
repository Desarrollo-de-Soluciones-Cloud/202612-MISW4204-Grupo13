package domain

import "errors"

var (
	ErrAssignmentNotFound          = errors.New("assignment not found")
	ErrAssignmentInvalidInput      = errors.New("invalid input")
	ErrAssignmentRoleRequired      = errors.New("assignment role is required")
	ErrAssignmentRoleInvalid       = errors.New("assignment role is invalid")
	ErrAssignmentWeeklyHoursInvalid = errors.New("assignment weekly hours must be greater than 0")
	ErrAssignmentUserIDRequired    = errors.New("assignment user id is required")
	ErrAssignmentWorkspaceIDRequired = errors.New("assignment workspace id is required")
)
