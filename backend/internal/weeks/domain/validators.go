package domain

import "time"

const (
	MinWeekNumber = 1
	MaxWeekNumber = 52
)

func ValidateWeekNumber(number int) error {
	switch {
	case number == 0:
		return ErrWeekNumberRequired
	case number < MinWeekNumber || number > MaxWeekNumber:
		return ErrWeekNumberInvalid
	default:
		return nil
	}
}

func ValidateWeekInitialDate(date time.Time) error {
	if date.IsZero() {
		return ErrWeekInitialDateRequired
	}
	return nil
}

func ValidateWeekFinalDate(date time.Time) error {
	if date.IsZero() {
		return ErrWeekFinalDateRequired
	}
	return nil
}

func ValidateWeekDateSequence(initialDate, finalDate time.Time) error {
	if initialDate.After(finalDate) {
		return ErrWeekDateSequenceInvalid
	}
	return nil
}

func ValidateWeekPeriodID(periodID uint) error {
	if periodID == 0 {
		return ErrWeekPeriodIDRequired
	}
	return nil
}
