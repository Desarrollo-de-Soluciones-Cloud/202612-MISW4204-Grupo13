package delivery

import (
	"backend/internal/periods/domain"
	"encoding/json"
	"errors"
	"strings"

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
		case "Name":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrPeriodNameRequired)
			}
		case "InitialDate":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrPeriodInitialDateRequired)
			}
		case "WeeksCount":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrPeriodWeeksCountRequired)
			}
		case "PeriodState":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrPeriodStateRequired)
			}
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrInvalidInput}
	}

	return result
}

// ParseAndValidateCreatePeriodRequest extracts and validates the CreatePeriodRequest field by field
// Returns the request object if valid, or a slice of accumulated errors
func ParseAndValidateCreatePeriodRequest(body []byte) (*CreatePeriodRequest, []error) {
	var rawData map[string]interface{}
	var validationErrors []error

	// Parse JSON
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, []error{domain.ErrInvalidInput}
	}

	req := &CreatePeriodRequest{}

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

	// Extract and validate InitialDate
	if dateVal, ok := rawData["initial_date"]; ok {
		if dateStr, ok := dateVal.(string); ok {
			req.InitialDate = dateStr
			// Validate initial date
			if strings.TrimSpace(req.InitialDate) == "" {
				validationErrors = append(validationErrors, domain.ErrPeriodInitialDateRequired)
			} else {
				if err := domain.ValidatePeriodInitialDate(req.InitialDate); err != nil {
					validationErrors = append(validationErrors, err)
				} else {
					// Only validate these if the date format is correct
					if err := domain.ValidatePeriodInitialDateIsMonday(req.InitialDate); err != nil {
						validationErrors = append(validationErrors, err)
					}
					if err := domain.ValidatePeriodInitialDateIsFuture(req.InitialDate); err != nil {
						validationErrors = append(validationErrors, err)
					}
				}
			}
		} else {
			validationErrors = append(validationErrors, domain.ErrPeriodInitialDateRequired)
		}
	} else {
		validationErrors = append(validationErrors, domain.ErrPeriodInitialDateRequired)
	}

	// Extract and validate WeeksCount
	if weeksVal, ok := rawData["weeks_count"]; ok {
		if weeksFloat, ok := weeksVal.(float64); ok {
			weeksInt := int(weeksFloat)
			req.WeeksCount = &weeksInt
			// Validate weeks count
			if err := domain.ValidatePeriodWeeksCount(weeksInt); err != nil {
				validationErrors = append(validationErrors, err)
			}
		} else {
			validationErrors = append(validationErrors, domain.ErrPeriodWeeksCountRequired)
		}
	} else {
		validationErrors = append(validationErrors, domain.ErrPeriodWeeksCountRequired)
	}

	// Extract and validate PeriodState
	if stateVal, ok := rawData["period_state"]; ok {
		if stateStr, ok := stateVal.(string); ok {
			req.PeriodState = stateStr
			// Validate state
			if strings.TrimSpace(req.PeriodState) == "" {
				validationErrors = append(validationErrors, domain.ErrPeriodStateRequired)
			} else {
				if err := domain.ValidatePeriodState(domain.PeriodState(req.PeriodState)); err != nil {
					validationErrors = append(validationErrors, err)
				}
			}
		} else {
			validationErrors = append(validationErrors, domain.ErrPeriodStateRequired)
		}
	} else {
		validationErrors = append(validationErrors, domain.ErrPeriodStateRequired)
	}

	// Return request and errors (if any)
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	return req, nil
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
