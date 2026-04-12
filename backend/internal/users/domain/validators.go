package domain

import (
	"regexp"
	"strings"
)

const (
	minUserNameLength  = 3
	maxUserNameLength  = 100
	minPasswordLength  = 8
	maxPasswordLength  = 72
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateUserName(name string) error {
	trimmedName := strings.TrimSpace(name)

	switch {
	case trimmedName == "":
		return ErrUserNameRequired
	case len(trimmedName) < minUserNameLength:
		return ErrUserNameTooShort
	case len(trimmedName) > maxUserNameLength:
		return ErrUserNameTooLong
	default:
		return nil
	}
}

func ValidateUserEmail(email string) error {
	normalizedEmail := NormalizeEmail(email)

	switch {
	case normalizedEmail == "":
		return ErrUserEmailRequired
	case !emailRegex.MatchString(normalizedEmail):
		return ErrUserEmailInvalid
	default:
		return nil
	}
}

func ValidateUserPassword(password string) error {
	switch {
	case strings.TrimSpace(password) == "":
		return ErrUserPasswordRequired
	case len(password) < minPasswordLength:
		return ErrUserPasswordTooShort
	case len(password) > maxPasswordLength:
		return ErrUserPasswordTooLong
	default:
		return nil
	}
}

func ValidateUserRole(role UserRole) error {
	switch {
	case strings.TrimSpace(string(role)) == "":
		return ErrUserRoleRequired
	case !IsValidUserRole(role):
		return ErrUserRoleInvalid
	default:
		return nil
	}
}
