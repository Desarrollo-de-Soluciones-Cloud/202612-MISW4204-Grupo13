package delivery

import (
	"backend/internal/tasks/domain"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

const dateLayout = "2006-01-02"

var taskBindingErrors = map[string]map[string]error{
	"AssignmentID": {
		"required": domain.ErrTaskAssignmentIDRequired,
	},
	"Title": {
		"required": domain.ErrTaskTitleRequired,
	},
	"Description": {
		"required": domain.ErrTaskDescriptionRequired,
	},
	"Status": {
		"required": domain.ErrTaskStatusRequired,
	},
	"SpentHours": {
		"required": domain.ErrTaskSpentHoursRequired,
	},
	"WeekStartDate": {
		"required": domain.ErrTaskWeekStartDateRequired,
	},
}

var taskUnmarshalTypeErrors = map[string]error{
	"assignment_id":  domain.ErrTaskAssignmentIDRequired,
	"AssignmentID":   domain.ErrTaskAssignmentIDRequired,
	"spent_hours":    domain.ErrTaskSpentHoursInvalid,
	"SpentHours":     domain.ErrTaskSpentHoursInvalid,
	"week_start_date": domain.ErrTaskWeekStartDateInvalid,
	"WeekStartDate":  domain.ErrTaskWeekStartDateInvalid,
}

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		result := make([]error, 0, len(validationErrors))
		for _, validationError := range validationErrors {
			if mappedErr := mapTaskBindingError(validationError); mappedErr != nil {
				result = append(result, mappedErr)
			}
		}

		if len(result) > 0 {
			return result
		}
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		if mappedErr := taskUnmarshalTypeErrors[unmarshalTypeError.Field]; mappedErr != nil {
			return []error{mappedErr}
		}
		return []error{domain.ErrInvalidInput}
	}

	return []error{domain.ErrInvalidInput}
}

func mapTaskBindingError(validationError validator.FieldError) error {
	fieldErrors, ok := taskBindingErrors[validationError.StructField()]
	if !ok {
		return nil
	}

	return fieldErrors[validationError.Tag()]
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
