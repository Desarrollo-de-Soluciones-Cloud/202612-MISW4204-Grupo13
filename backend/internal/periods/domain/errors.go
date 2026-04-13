package domain

import "errors"

var (
	ErrPeriodNotFound                        = errors.New("period not found")
	ErrInvalidInput                          = errors.New("invalid input")
	ErrPeriodNameAlreadyExists               = errors.New("a period with the same name already exists")
	ErrPeriodNameRequired                    = errors.New("period name is required")
	ErrPeriodNameWrongFormat                 = errors.New("period name must have format YYYY-SS. For instance, 2026-10")
	ErrPeriodInitialDateRequired             = errors.New("period initial date is required")
	ErrPeriodInitialDateWrongFormat          = errors.New("period initial date must have format YYYY-MM-DD. For instance, 2026-10-01")
	ErrPeriodInitialDateMustBeMonday         = errors.New("period initial date must be a monday")
	ErrPeriodInitialDateMustBeFuture         = errors.New("period initial date must be after current date")
	ErrPeriodWeeksCountRequired              = errors.New("period weeks count is required")
	ErrPeriodWeeksCountInvalid               = errors.New("period weeks count is invalid. Allowed values: 8, 16")
	ErrPeriodStateRequired                   = errors.New("period state is required")
	ErrPeriodStateInvalid                    = errors.New("period state is invalid. The following states are allowed: " + ValidPeriodStatesString())
	ErrPeriodCannotBeUpdated                 = errors.New("period cannot be updated after its initial date has passed")
)
