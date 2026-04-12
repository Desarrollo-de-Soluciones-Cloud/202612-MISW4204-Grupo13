package delivery

import (
	"backend/internal/auth/domain"
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
		case "Email":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrAuthEmailRequired)
			case "email":
				result = append(result, domain.ErrAuthEmailInvalid)
			}
		case "Password":
			switch validationError.Tag() {
			case "required":
				result = append(result, domain.ErrAuthPasswordRequired)
			case "min":
				result = append(result, domain.ErrAuthPasswordTooShort)
			case "max":
				result = append(result, domain.ErrAuthPasswordTooLong)
			}
		}
	}

	if len(result) == 0 {
		return []error{domain.ErrInvalidInput}
	}

	return result
}

func extractBearerToken(authorizationHeader string) (string, error) {
	trimmedHeader := strings.TrimSpace(authorizationHeader)
	if trimmedHeader == "" {
		return "", domain.ErrAuthTokenRequired
	}

	parts := strings.SplitN(trimmedHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", domain.ErrAuthTokenInvalid
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", domain.ErrAuthTokenRequired
	}

	return token, nil
}
