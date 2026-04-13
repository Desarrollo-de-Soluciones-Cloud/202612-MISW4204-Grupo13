package domain

import "errors"

var (
	ErrWeekNotFound                  = errors.New("week not found")
	ErrInvalidInput                  = errors.New("invalid input")
	ErrWeekPeriodIDRequired          = errors.New("week period id is required")
	ErrWeekNumberInvalid             = errors.New("week number is invalid")
	ErrWeekInitialDateRequired       = errors.New("week initial date is required")
	ErrWeekInitialDateWrongFormat    = errors.New("week initial date must have format YYYY-MM-DD. For instance, 2026-10-05")
	ErrWeekInitialDateMustBeMonday   = errors.New("week initial date must be a monday")
	ErrWeekFinalDateRequired         = errors.New("week final date is required")
	ErrWeekFinalDateWrongFormat      = errors.New("week final date must have format YYYY-MM-DD. For instance, 2026-10-11")
	ErrWeekFinalDateMustBeSunday     = errors.New("week final date must be a sunday")
	ErrWeekDateRangeInvalid          = errors.New("week date range is invalid. The final date must be exactly 6 days after the initial date.")
	ErrWeekFinalDateMismatch         = errors.New("week final date does not match the expected final date for the given initial date and week count")
	ErrWeekCountInvalid              = errors.New("week count is invalid. Allowed values: 8, 16")
	ErrWeeksAlreadyExistForPeriod    = errors.New("weeks already exist for the given period")
)
