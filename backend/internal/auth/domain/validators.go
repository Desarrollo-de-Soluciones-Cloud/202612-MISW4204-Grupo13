package domain

import (
	"regexp"
	"strings"
)

const (
	minAuthPasswordLength = 8
	maxAuthPasswordLength = 72
)

var authEmailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func NormalizeAuthEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateAuthEmail(email string) error {
	normalizedEmail := NormalizeAuthEmail(email)

	switch {
	case normalizedEmail == "":
		return ErrAuthEmailRequired
	case !authEmailRegex.MatchString(normalizedEmail):
		return ErrAuthEmailInvalid
	default:
		return nil
	}
}

func ValidateAuthPassword(password string) error {
	switch {
	case strings.TrimSpace(password) == "":
		return ErrAuthPasswordRequired
	case len(password) < minAuthPasswordLength:
		return ErrAuthPasswordTooShort
	case len(password) > maxAuthPasswordLength:
		return ErrAuthPasswordTooLong
	default:
		return nil
	}
}

func ValidateTokenString(token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrAuthTokenRequired
	}

	return nil
}
