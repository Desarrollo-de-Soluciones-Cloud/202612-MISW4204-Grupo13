package domain

import "errors"

var (
	ErrWeekNumberRequired        = errors.New("week number is required")
	ErrWeekNumberInvalid         = errors.New("week number must be between 1 and 52")
	ErrWeekInitialDateRequired   = errors.New("week initial date is required")
	ErrWeekFinalDateRequired     = errors.New("week final date is required")
	ErrWeekDateSequenceInvalid   = errors.New("week initial date must be before or equal to final date")
	ErrWeekPeriodIDRequired      = errors.New("week period id is required")
	ErrWeekNotFound              = errors.New("week not found")
	ErrWeekAlreadyExists         = errors.New("week already exists")
)
