package domain

import "errors"

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrAuthEmailRequired    = errors.New("auth email is required")
	ErrAuthEmailInvalid     = errors.New("auth email is invalid")
	ErrAuthPasswordRequired = errors.New("auth password is required")
	ErrAuthPasswordTooShort = errors.New("auth password must have at least 8 characters")
	ErrAuthPasswordTooLong  = errors.New("auth password must have at most 72 characters")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAuthTokenRequired    = errors.New("authentication token is required")
	ErrAuthTokenInvalid     = errors.New("authentication token is invalid")
	ErrAuthTokenExpired     = errors.New("authentication token has expired")
	ErrAuthForbidden        = errors.New("insufficient permissions")
)
