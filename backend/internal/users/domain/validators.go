package domain

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	minUserNameLength  = 3
	maxUserNameLength  = 100
	minPasswordLength  = 8
	maxPasswordLength  = 72
	passwordHashCost   = bcrypt.DefaultCost
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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return "", ErrPasswordHashingFailed
	}

	return string(hash), nil
}
