package delivery

import (
	"backend/internal/users/domain"
	"errors"

	"github.com/go-playground/validator/v10"
)

var userBindingErrors = map[string]map[string]error{
	"Name": {
		"required": domain.ErrUserNameRequired,
		"min":      domain.ErrUserNameTooShort,
		"max":      domain.ErrUserNameTooLong,
	},
	"Email": {
		"required": domain.ErrUserEmailRequired,
		"email":    domain.ErrUserEmailInvalid,
	},
	"Password": {
		"required": domain.ErrUserPasswordRequired,
		"min":      domain.ErrUserPasswordTooShort,
		"max":      domain.ErrUserPasswordTooLong,
	},
	"GlobalRole": {
		"required": domain.ErrUserRoleRequired,
	},
}

func mapBindingErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []error{domain.ErrInvalidInput}
	}

	result := make([]error, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		if mappedErr := mapUserBindingError(validationError); mappedErr != nil {
			result = append(result, mappedErr)
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrInvalidInput}
	}

	return result
}

func mapUserBindingError(validationError validator.FieldError) error {
	fieldErrors, ok := userBindingErrors[validationError.StructField()]
	if !ok {
		return nil
	}

	return fieldErrors[validationError.Tag()]
}
