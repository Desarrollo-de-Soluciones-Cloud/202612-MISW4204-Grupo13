package delivery

import (
	"backend/internal/workspaces/domain"
	"errors"

	"github.com/go-playground/validator/v10"
)

var workspaceBindingErrors = map[string]map[string]error{
	"PeriodID": {
		"required": domain.ErrWorkspacePeriodIDRequired,
	},
	"UserID": {
		"required": domain.ErrWorkspaceUserIDRequired,
	},
	"Name": {
		"required": domain.ErrWorkspaceNameRequired,
	},
	"Type": {
		"required": domain.ErrWorkspaceTypeRequired,
		"oneof":    domain.ErrWorkspaceTypeInvalid,
	},
	"InitialDate": {
		"required": domain.ErrWorkspaceInitialDateRequired,
	},
	"FinalDate": {
		"required": domain.ErrWorkspaceFinalDateRequired,
	},
	"Observations": {
		"required": errors.New("observations is required"),
	},
	"State": {
		"required": domain.ErrWorkspaceStateRequired,
		"oneof":    domain.ErrWorkspaceStateInvalid,
	},
}

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []error{domain.ErrInvalidInput}
	}

	result := make([]error, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		if mappedErr := mapWorkspaceBindingError(validationError); mappedErr != nil {
			result = append(result, mappedErr)
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrInvalidInput}
	}

	return result
}

func mapWorkspaceBindingError(validationError validator.FieldError) error {
	fieldErrors, ok := workspaceBindingErrors[validationError.StructField()]
	if !ok {
		return nil
	}

	return fieldErrors[validationError.Tag()]
}

func isWorkspaceValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrWorkspacePeriodNotFound) ||
		errors.Is(err, domain.ErrWorkspacePeriodClosed) ||
		errors.Is(err, domain.ErrWorkspaceInscriptionClosed) ||
		errors.Is(err, domain.ErrWorkspaceUserNotFound) ||
		errors.Is(err, domain.ErrWorkspaceUserNotProfessor) ||
		errors.Is(err, domain.ErrWorkspaceClosedUpdateForbidden) ||
		errors.Is(err, domain.ErrWorkspaceUserIDChangeNotAllowed) ||
		errors.Is(err, domain.ErrWorkspaceNameRequired) ||
		errors.Is(err, domain.ErrWorkspaceNameTooLong) ||
		errors.Is(err, domain.ErrWorkspacePeriodIDRequired) ||
		errors.Is(err, domain.ErrWorkspaceUserIDRequired) ||
		errors.Is(err, domain.ErrWorkspaceTypeRequired) ||
		errors.Is(err, domain.ErrWorkspaceTypeInvalid) ||
		errors.Is(err, domain.ErrWorkspaceInitialDateRequired) ||
		errors.Is(err, domain.ErrWorkspaceInitialDateWrongFormat) ||
		errors.Is(err, domain.ErrWorkspaceInitialDateOutOfRange) ||
		errors.Is(err, domain.ErrWorkspaceFinalDateRequired) ||
		errors.Is(err, domain.ErrWorkspaceFinalDateWrongFormat) ||
		errors.Is(err, domain.ErrWorkspaceFinalDateOutOfRange) ||
		errors.Is(err, domain.ErrWorkspaceDateSequenceInvalid) ||
		errors.Is(err, domain.ErrWorkspaceStateRequired) ||
		errors.Is(err, domain.ErrWorkspaceStateInvalid)
}
