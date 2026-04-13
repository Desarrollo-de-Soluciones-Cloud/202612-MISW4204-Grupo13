package domain

import "errors"

var (
	// Workspace Name errors
	ErrWorkspaceNameRequired = errors.New("workspace name is required")
	ErrWorkspaceNameTooLong  = errors.New("workspace name is too long")

	// Workspace IDs errors
	ErrWorkspacePeriodIDRequired = errors.New("workspace period_id is required")
	ErrWorkspaceUserIDRequired   = errors.New("workspace user_id is required")

	// Workspace Type errors
	ErrWorkspaceTypeRequired = errors.New("workspace type is required")

	// Workspace Date errors
	ErrWorkspaceInitialDateRequired   = errors.New("workspace initial_date is required")
	ErrWorkspaceInitialDateWrongFormat  = errors.New("workspace initial_date has wrong format (expected YYYY-MM-DD)")
	ErrWorkspaceFinalDateRequired      = errors.New("workspace final_date is required")
	ErrWorkspaceFinalDateWrongFormat   = errors.New("workspace final_date has wrong format (expected YYYY-MM-DD)")
	ErrWorkspaceDateSequenceInvalid    = errors.New("workspace initial_date must be before final_date")

	// Workspace State errors
	ErrWorkspaceStateRequired = errors.New("workspace state is required")
	ErrWorkspaceStateInvalid  = errors.New("workspace state is invalid")
)
