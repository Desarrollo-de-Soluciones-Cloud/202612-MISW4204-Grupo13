package delivery

import (
	"backend/internal/users/domain"
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
				result = append(result, domain.ErrUserNameRequired)
			case "min":
				result = append(result, domain.ErrUserNameTooShort)
			case "max":
				result = append(result, domain.ErrUserNameTooLong)
			}
		case "Email":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrUserEmailRequired)
			case "email":
				result = append(result, domain.ErrUserEmailInvalid)
			}
		case "Password":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrUserPasswordRequired)
			case "min":
				result = append(result, domain.ErrUserPasswordTooShort)
			case "max":
				result = append(result, domain.ErrUserPasswordTooLong)
			}
		case "GlobalRole":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrUserRoleRequired)
			}
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrInvalidInput}
	}

	return result
}
