package domain

import "strings"

func ValidateAssignmentUserID(userID uint) error {
	if userID == 0 {
		return ErrAssignmentUserIDRequired
	}

	return nil
}

func ValidateAssignmentWorkspaceID(workspaceID uint) error {
	if workspaceID == 0 {
		return ErrAssignmentWorkspaceIDRequired
	}

	return nil
}

func ValidateAssignmentRole(role AssignmentRole) error {
	switch {
	case strings.TrimSpace(string(role)) == "":
		return ErrAssignmentRoleRequired
	case !IsValidAssignmentRole(role):
		return ErrAssignmentRoleInvalid
	default:
		return nil
	}
}

func ValidateAssignmentWeeklyHours(weeklyHours int) error {
	if weeklyHours <= 0 {
		return ErrAssignmentWeeklyHoursInvalid
	}

	return nil
}
