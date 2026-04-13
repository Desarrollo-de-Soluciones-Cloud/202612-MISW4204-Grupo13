package delivery

import (
	"backend/internal/tasks/domain"
	"encoding/json"
	"errors"
	"time"

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
			case "WeekStartDate":
				if validationError.Tag() == "required" {
					result = append(result, domain.ErrTaskWeekStartDateRequired)
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
		case "spent_hours", "SpentHours":
			return []error{domain.ErrTaskSpentHoursInvalid}
		case "week_start_date", "WeekStartDate":
			return []error{domain.ErrTaskWeekStartDateInvalid}
		default:
			return []error{domain.ErrInvalidInput}
		}
	}

	return []error{domain.ErrInvalidInput}
}

func parseWeekStartDate(rawDate string) (time.Time, error) {
	if rawDate == "" {
		return time.Time{}, domain.ErrTaskWeekStartDateRequired
	}

	parsedDate, err := time.Parse(dateLayout, rawDate)
	if err != nil {
		return time.Time{}, domain.ErrTaskWeekStartDateInvalid
	}

	return parsedDate, nil
}
