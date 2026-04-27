package delivery

import (
	reportsDomain "backend/internal/reports/domain"
	sharedHelpers "backend/internal/shared/helpers"
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
)

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			return []error{reportsDomain.ErrReportInvalidInput}
		}
		return []error{reportsDomain.ErrReportInvalidInput}
	}

	result := make([]error, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		switch validationError.StructField() {
		case "WorkspaceID":
			switch validationError.Tag() {
			case "required":
				result = append(result, reportsDomain.ErrReportWorkspaceIDRequired)
			}
		case "WeekID":
			switch validationError.Tag() {
			case "required":
				result = append(result, reportsDomain.ErrReportWeekIDRequired)
			}
		}
	}

	if len(result) == 0 {
		return []error{reportsDomain.ErrReportInvalidInput}
	}

	return result
}

func isReportValidationError(err error) bool {
	return errors.Is(err, reportsDomain.ErrReportInvalidInput) ||
		errors.Is(err, reportsDomain.ErrReportWorkspaceIDRequired) ||
		errors.Is(err, reportsDomain.ErrReportWeekIDRequired) ||
		errors.Is(err, reportsDomain.ErrReportAssignmentIDRequired) ||
		errors.Is(err, reportsDomain.ErrReportUserIDRequired) ||
		errors.Is(err, reportsDomain.ErrReportFilePathRequired) ||
		errors.Is(err, reportsDomain.ErrReportWorkspaceFilterRequired)
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