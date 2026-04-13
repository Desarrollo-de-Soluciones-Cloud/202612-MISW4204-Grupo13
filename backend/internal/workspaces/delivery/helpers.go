package delivery

import (
	"backend/internal/workspaces/domain"
	"errors"

	"github.com/go-playground/validator/v10"
)

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []error{domain.ErrInvalidInput}
	}

	result := make([]error, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		switch validationError.StructField() {
		case "PeriodID":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrWorkspacePeriodIDRequired)
			}
		case "UserID":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrWorkspaceUserIDRequired)
			}
		case "Name":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrWorkspaceNameRequired)
			}
		case "Type":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrWorkspaceTypeRequired)
			}
		case "InitialDate":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrWorkspaceInitialDateRequired)
			}
		case "FinalDate":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrWorkspaceFinalDateRequired)
			}
		case "Observations":
			switch validationError.Tag() {
			case "required":
				result = append(result, errors.New("observations is required"))
			}
		case "State":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrWorkspaceStateRequired)
			}
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrInvalidInput}
	}

	return result
}

func isWorkspaceValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrWorkspacePeriodNotFound) ||
		errors.Is(err, domain.ErrWorkspaceUserNotFound) ||
		errors.Is(err, domain.ErrWorkspaceUserNotProfessor) ||
		errors.Is(err, domain.ErrWorkspaceNameRequired) ||
		errors.Is(err, domain.ErrWorkspaceNameTooLong) ||
		errors.Is(err, domain.ErrWorkspacePeriodIDRequired) ||
		errors.Is(err, domain.ErrWorkspaceUserIDRequired) ||
		errors.Is(err, domain.ErrWorkspaceTypeRequired) ||
		errors.Is(err, domain.ErrWorkspaceInitialDateRequired) ||
		errors.Is(err, domain.ErrWorkspaceInitialDateWrongFormat) ||
		errors.Is(err, domain.ErrWorkspaceFinalDateRequired) ||
		errors.Is(err, domain.ErrWorkspaceFinalDateWrongFormat) ||
		errors.Is(err, domain.ErrWorkspaceDateSequenceInvalid) ||
		errors.Is(err, domain.ErrWorkspaceStateRequired) ||
		errors.Is(err, domain.ErrWorkspaceStateInvalid)
}
