package domain

import (
	"strings"
	"time"
)

func NormalizeWorkspaceName(name string) string {
	return strings.TrimSpace(name)
}

func ValidateWorkspaceName(name string) error {
	trimmedName := strings.TrimSpace(name)

	switch {
	case trimmedName == "":
		return ErrWorkspaceNameRequired
	case len(trimmedName) > 100:
		return ErrWorkspaceNameTooLong
	default:
		return nil
	}
}

func ValidateWorkspacePeriodID(periodID uint) error {
	if periodID == 0 {
		return ErrWorkspacePeriodIDRequired
	}
	return nil
}

func ValidateWorkspaceUserID(userID uint) error {
	if userID == 0 {
		return ErrWorkspaceUserIDRequired
	}
	return nil
}

func ValidateWorkspaceType(workspaceType WorkspaceType) error {
	trimmedType := strings.TrimSpace(string(workspaceType))

	if trimmedType == "" {
		return ErrWorkspaceTypeRequired
	}
	// Aquí puedes añadir más validaciones si hay tipos específicos válidos
	return nil
}

func ValidateWorkspaceInitialDate(date string) error {
	trimmedDate := strings.TrimSpace(date)

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrWorkspaceInitialDateWrongFormat
	}

	switch {
	case trimmedDate == "":
		return ErrWorkspaceInitialDateRequired
	case len(trimmedDate) != 10:
		return ErrWorkspaceInitialDateWrongFormat
	default:
		return nil
	}
}

func ValidateWorkspaceFinalDate(date string) error {
	trimmedDate := strings.TrimSpace(date)

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrWorkspaceFinalDateWrongFormat
	}

	switch {
	case trimmedDate == "":
		return ErrWorkspaceFinalDateRequired
	case len(trimmedDate) != 10:
		return ErrWorkspaceFinalDateWrongFormat
	default:
		return nil
	}
}

func ValidateWorkspaceDateSequence(initialDate, finalDate string) error {
	if initialDate > finalDate {
		return ErrWorkspaceDateSequenceInvalid
	}
	return nil
}

func ValidateWorkspaceState(state WorkspaceState) error {
	switch {
	case strings.TrimSpace(string(state)) == "":
		return ErrWorkspaceStateRequired
	case !IsValidWorkspaceState(state):
		return ErrWorkspaceStateInvalid
	default:
		return nil
	}
}
