package delivery

import (
	"backend/internal/users/domain"
	"errors"

	"github.com/go-playground/validator/v10"
)

func mapBindingError(err error) error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return domain.ErrInvalidInput
	}

	validationError := validationErrors[0]

	switch validationError.StructField() {
	case "Name":
		switch validationError.Tag() {
		case "required":
			return domain.ErrUserNameRequired
		case "min":
			return domain.ErrUserNameTooShort
		case "max":
			return domain.ErrUserNameTooLong
		}
	case "Email":
		switch validationError.Tag() {
		case "required":
			return domain.ErrUserEmailRequired
		case "email":
			return domain.ErrUserEmailInvalid
		}
	case "Password":
		switch validationError.Tag() {
		case "required":
			return domain.ErrUserPasswordRequired
		case "min":
			return domain.ErrUserPasswordTooShort
		case "max":
			return domain.ErrUserPasswordTooLong
		}
	case "GlobalRole":
		switch validationError.Tag() {
		case "required":
			return domain.ErrUserRoleRequired
		}
	}

	return domain.ErrInvalidInput
}
