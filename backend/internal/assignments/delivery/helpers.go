package delivery

import (
	"backend/internal/assignments/domain"
	"errors"

	"github.com/go-playground/validator/v10"
)

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []error{domain.ErrAssignmentInvalidInput}
	}

	result := make([]error, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		switch validationError.StructField() {
		case "UserID":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrAssignmentUserIDRequired)
			}
		case "WorkspaceID":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrAssignmentWorkspaceIDRequired)
			}
		case "Role":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrAssignmentRoleRequired)
			}
		case "WeeklyHours":
			switch validationError.Tag() {
			case "required", "min":
				result = append(result, domain.ErrAssignmentWeeklyHoursInvalid)
			}
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrAssignmentInvalidInput}
	}

	return result
}

func isAssignmentValidationError(err error) bool {
	return errors.Is(err, domain.ErrAssignmentInvalidInput) ||
		errors.Is(err, domain.ErrAssignmentRoleRequired) ||
		errors.Is(err, domain.ErrAssignmentRoleInvalid) ||
		errors.Is(err, domain.ErrAssignmentWeeklyHoursInvalid) ||
		errors.Is(err, domain.ErrAssignmentUserIDRequired) ||
		errors.Is(err, domain.ErrAssignmentWorkspaceIDRequired)
}
