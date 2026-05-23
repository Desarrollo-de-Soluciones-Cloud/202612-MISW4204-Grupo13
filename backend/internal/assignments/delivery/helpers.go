package delivery

import (
	"backend/internal/assignments/domain"
	"errors"

	"github.com/go-playground/validator/v10"
)

var assignmentBindingErrors = map[string]map[string]error{
	"UserID": {
		"required": domain.ErrAssignmentUserIDRequired,
	},
	"WorkspaceID": {
		"required": domain.ErrAssignmentWorkspaceIDRequired,
	},
	"Role": {
		"required": domain.ErrAssignmentRoleRequired,
	},
	"WeeklyHours": {
		"required": domain.ErrAssignmentWeeklyHoursInvalid,
		"min":      domain.ErrAssignmentWeeklyHoursInvalid,
	},
}

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []error{domain.ErrAssignmentInvalidInput}
	}

	result := make([]error, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		if mappedErr := mapAssignmentBindingError(validationError); mappedErr != nil {
			result = append(result, mappedErr)
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrAssignmentInvalidInput}
	}

	return result
}

func mapAssignmentBindingError(validationError validator.FieldError) error {
	fieldErrors, ok := assignmentBindingErrors[validationError.StructField()]
	if !ok {
		return nil
	}

	return fieldErrors[validationError.Tag()]
}

func isAssignmentValidationError(err error) bool {
	return errors.Is(err, domain.ErrAssignmentInvalidInput) ||
		errors.Is(err, domain.ErrAssignmentRoleRequired) ||
		errors.Is(err, domain.ErrAssignmentRoleInvalid) ||
		errors.Is(err, domain.ErrAssignmentWeeklyHoursInvalid) ||
		errors.Is(err, domain.ErrAssignmentUserIDRequired) ||
		errors.Is(err, domain.ErrAssignmentWorkspaceIDRequired) ||
		errors.Is(err, domain.ErrAssignmentUserInvalidRole) ||
		errors.Is(err, domain.ErrAssignmentRoleNotAllowedForUser) ||
		errors.Is(err, domain.ErrAssignmentProfessorCannotChangeWeeklyHours) ||
		errors.Is(err, domain.ErrAssignmentAssistantHoursLimitExceeded) ||
		errors.Is(err, domain.ErrAssignmentMonitorCountLimitExceeded) ||
		errors.Is(err, domain.ErrAssignmentMonitorHoursLimitExceeded) ||
		errors.Is(err, domain.ErrAssignmentMonitorFortyPercentExceeded)
}
