package domain

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserEmailAlreadyInUse = errors.New("user email is already in use")
	ErrInvalidInput          = errors.New("invalid input")
	ErrUserNameRequired      = errors.New("user name is required")
	ErrUserNameTooShort      = errors.New("user name must have at least 3 characters")
	ErrUserNameTooLong       = errors.New("user name must have at most 100 characters")
	ErrUserEmailRequired     = errors.New("user email is required")
	ErrUserEmailInvalid      = errors.New("user email is invalid")
	ErrUserPasswordRequired  = errors.New("user password is required")
	ErrUserPasswordTooShort  = errors.New("user password must have at least 8 characters")
	ErrUserPasswordTooLong   = errors.New("user password must have at most 72 characters")
	ErrUserRoleRequired      = errors.New("user global role is required")
	ErrUserRoleInvalid       = errors.New("user global role is invalid")
)
