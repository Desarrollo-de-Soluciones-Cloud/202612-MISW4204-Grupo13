package delivery

import (
	"backend/internal/tasks/domain"
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
)

const dateLayout = "2006-01-02"

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		result := make([]error, 0, len(validationErrors))
		for _, validationError := range validationErrors {
			switch validationError.StructField() {
			case "AssignmentID":
				if validationError.Tag() == "required" {
					result = append(result, domain.ErrTaskAssignmentIDRequired)
				}
			case "WeekID":
				if validationError.Tag() == "required" {
					result = append(result, domain.ErrTaskWeekIDRequired)
				}
			case "Title":
				if validationError.Tag() == "required" {
					result = append(result, domain.ErrTaskTitleRequired)
				}
			case "Description":
				if validationError.Tag() == "required" {
					result = append(result, domain.ErrTaskDescriptionRequired)
				}
			case "Status":
				if validationError.Tag() == "required" {
					result = append(result, domain.ErrTaskStatusRequired)
				}
			case "SpentHours":
				if validationError.Tag() == "required" {
					result = append(result, domain.ErrTaskSpentHoursRequired)
				}
			}
		}

		if len(result) > 0 {
			return result
		}
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		switch unmarshalTypeError.Field {
		case "assignment_id", "AssignmentID":
			return []error{domain.ErrTaskAssignmentIDRequired}
		case "week_id", "WeekID":
			return []error{domain.ErrTaskWeekIDRequired}
		case "spent_hours", "SpentHours":
			return []error{domain.ErrTaskSpentHoursInvalid}
		default:
			return []error{domain.ErrInvalidInput}
		}
	}

	return []error{domain.ErrInvalidInput}
}
