package delivery

import (
	reportsDomain "backend/internal/reports/domain"
	sharedHelpers "backend/internal/shared/helpers"
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
)

func mapBindingErrors(err error) []error {
	if mapped := mapValidationBindingErrors(err); len(mapped) > 0 {
		return mapped
	}

	if mapped := mapUnmarshalBindingError(err); len(mapped) > 0 {
		return mapped
	}

	return []error{reportsDomain.ErrReportInvalidInput}
}

func mapValidationBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		result := make([]error, 0, len(validationErrors))

		errorByFieldAndTag := map[string]error{
			"WorkspaceID:required": reportsDomain.ErrReportWorkspaceIDRequired,
			"WeekID:required":      reportsDomain.ErrReportWeekIDRequired,
		}

		for _, validationError := range validationErrors {
			key := validationError.StructField() + ":" + validationError.Tag()
			if mapped, exists := errorByFieldAndTag[key]; exists {
				result = append(result, mapped)
			}
		}

		if len(result) > 0 {
			return result
		}
	}

	return nil
}

func mapUnmarshalBindingError(err error) []error {
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return []error{reportsDomain.ErrReportInvalidInput}
	}

	return nil
}

func parseRequiredWorkspaceID(raw string) (uint, error) {
	if raw == "" {
		return 0, reportsDomain.ErrReportWorkspaceFilterRequired
	}

	return sharedHelpers.ParseResourceID(raw)
}

func parseOptionalResourceID(raw string) (*uint, error) {
	if raw == "" {
		return nil, nil
	}

	id, err := sharedHelpers.ParseResourceID(raw)
	if err != nil {
		return nil, err
	}

	return &id, nil
}