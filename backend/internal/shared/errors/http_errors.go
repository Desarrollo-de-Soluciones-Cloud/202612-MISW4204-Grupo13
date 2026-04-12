package errors

import "errors"

var (
	ErrInvalidResourceID    = errors.New("invalid resource id")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("authentication required")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("resource conflict")
	ErrInternalServerError = errors.New("internal server error")
)
