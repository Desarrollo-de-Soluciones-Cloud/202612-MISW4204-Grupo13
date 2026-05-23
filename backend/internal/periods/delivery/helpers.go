package delivery

import (
	"backend/internal/periods/domain"
	"encoding/json"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var requiredFieldErrors = map[string]error{
	"Name":        domain.ErrPeriodNameRequired,
	"InitialDate": domain.ErrPeriodInitialDateRequired,
	"WeeksCount":  domain.ErrPeriodWeeksCountRequired,
	"PeriodState": domain.ErrPeriodStateRequired,
}

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []error{domain.ErrInvalidInput}
	}

	result := make([]error, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		if mappedErr := mapRequiredBindingError(validationError); mappedErr != nil {
			result = append(result, mappedErr)
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrInvalidInput}
	}

	return result
}

func mapRequiredBindingError(validationError validator.FieldError) error {
	if validationError.Tag() != "required" {
		return nil
	}

	return requiredFieldErrors[validationError.StructField()]
}

// ParseAndValidateCreatePeriodRequest extracts and validates the CreatePeriodRequest field by field
// Returns the request object if valid, or a slice of accumulated errors
func ParseAndValidateCreatePeriodRequest(body []byte) (*CreatePeriodRequest, []error) {
	var rawData map[string]interface{}

	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, []error{domain.ErrInvalidInput}
	}

	req := &CreatePeriodRequest{}
	validationErrors := make([]error, 0)

	validationErrors = append(validationErrors, validateCreatePeriodName(rawData, req)...)
	validationErrors = append(validationErrors, validateCreatePeriodInitialDate(rawData, req)...)
	validationErrors = append(validationErrors, validateCreatePeriodWeeksCount(rawData, req)...)
	validationErrors = append(validationErrors, validateCreatePeriodState(rawData, req)...)

	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return req, nil
}

func validateCreatePeriodName(rawData map[string]interface{}, req *CreatePeriodRequest) []error {
	name, ok := getRequiredString(rawData, "name")
	if !ok {
		return []error{domain.ErrPeriodNameRequired}
	}

	req.Name = name
	if err := domain.ValidatePeriodName(req.Name); err != nil {
		return []error{err}
	}

	return nil
}

func validateCreatePeriodInitialDate(rawData map[string]interface{}, req *CreatePeriodRequest) []error {
	initialDate, ok := getRequiredString(rawData, "initial_date")
	if !ok {
		return []error{domain.ErrPeriodInitialDateRequired}
	}

	req.InitialDate = initialDate
	if err := domain.ValidatePeriodInitialDate(req.InitialDate); err != nil {
		return []error{err}
	}

	validationErrors := make([]error, 0, 2)
	if err := domain.ValidatePeriodInitialDateIsMonday(req.InitialDate); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := domain.ValidatePeriodInitialDateIsFuture(req.InitialDate); err != nil {
		validationErrors = append(validationErrors, err)
	}

	return validationErrors
}

func validateCreatePeriodWeeksCount(rawData map[string]interface{}, req *CreatePeriodRequest) []error {
	weeksCount, ok := getRequiredInt(rawData, "weeks_count")
	if !ok {
		return []error{domain.ErrPeriodWeeksCountRequired}
	}

	req.WeeksCount = &weeksCount
	if err := domain.ValidatePeriodWeeksCount(weeksCount); err != nil {
		return []error{err}
	}

	return nil
}

func validateCreatePeriodState(rawData map[string]interface{}, req *CreatePeriodRequest) []error {
	periodState, ok := getRequiredString(rawData, "period_state")
	if !ok {
		return []error{domain.ErrPeriodStateRequired}
	}

	req.PeriodState = periodState
	if err := domain.ValidatePeriodState(domain.PeriodState(req.PeriodState)); err != nil {
		return []error{err}
	}

	return nil
}

func getRequiredString(rawData map[string]interface{}, key string) (string, bool) {
	value, ok := rawData[key]
	if !ok {
		return "", false
	}

	stringValue, ok := value.(string)
	if !ok || strings.TrimSpace(stringValue) == "" {
		return "", false
	}

	return stringValue, true
}

func getRequiredInt(rawData map[string]interface{}, key string) (int, bool) {
	value, ok := rawData[key]
	if !ok {
		return 0, false
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0, false
	}

	return int(floatValue), true
}

// ParseAndValidateUpdatePeriodRequest extracts and validates the UpdatePeriodRequest field by field
func ParseAndValidateUpdatePeriodRequest(body []byte) (*UpdatePeriodRequest, []error) {
	var rawData map[string]interface{}
	var validationErrors []error

	// Parse JSON
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, []error{domain.ErrInvalidInput}
	}

	req := &UpdatePeriodRequest{}

	// Extract and validate Name
	if nameVal, ok := rawData["name"]; ok {
		if nameStr, ok := nameVal.(string); ok {
			req.Name = nameStr
			// Validate name
			if strings.TrimSpace(req.Name) == "" {
				validationErrors = append(validationErrors, domain.ErrPeriodNameRequired)
			} else {
				if err := domain.ValidatePeriodName(req.Name); err != nil {
					validationErrors = append(validationErrors, err)
				}
			}
		} else {
			validationErrors = append(validationErrors, domain.ErrPeriodNameRequired)
		}
	} else {
		validationErrors = append(validationErrors, domain.ErrPeriodNameRequired)
	}

	// Return request and errors (if any)
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return req, nil
}
