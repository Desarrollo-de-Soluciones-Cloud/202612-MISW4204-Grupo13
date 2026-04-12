package delivery

import (
	"backend/internal/periods/domain"
	"errors"

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
		case "FinalDate":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrPeriodFinalDateRequired)
			}
		case "InscriptionFinalDate":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrPeriodInscriptionFinalDateRequired)
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
