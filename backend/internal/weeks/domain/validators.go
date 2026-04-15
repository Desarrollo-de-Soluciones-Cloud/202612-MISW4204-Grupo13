package domain

import (
	"strings"
	"time"
)

const weekDateLayout = "2006-01-02"

func ValidateWeekPeriodID(periodID uint) error {
	if periodID == 0 {
		return ErrWeekPeriodIDRequired
	}
	return nil
}

func ValidateWeekNumber(number int) error {
	if number <= 0 {
		return ErrWeekNumberInvalid
	}
	return nil
}

func ValidateWeekInitialDate(initialDate string) error {
	switch {
	case strings.TrimSpace(initialDate) == "":
		return ErrWeekInitialDateRequired
	case !isValidWeekDate(initialDate):
		return ErrWeekInitialDateWrongFormat
	default:
		return nil
	}
}

func ValidateWeekFinalDate(finalDate string) error {
	switch {
	case strings.TrimSpace(finalDate) == "":
		return ErrWeekFinalDateRequired
	case !isValidWeekDate(finalDate):
		return ErrWeekFinalDateWrongFormat
	default:
		return nil
	}
}

func ValidateWeekInitialDateIsMonday(initialDate string) error {
	date, err := parseWeekDate(initialDate)
	if err != nil {
		return ErrWeekInitialDateWrongFormat
	}
	if date.Weekday() != time.Monday {
		return ErrWeekInitialDateMustBeMonday
	}
	return nil
}

func ValidateWeekFinalDateIsSunday(finalDate string) error {
	date, err := parseWeekDate(finalDate)
	if err != nil {
		return ErrWeekFinalDateWrongFormat
	}
	if date.Weekday() != time.Sunday {
		return ErrWeekFinalDateMustBeSunday
	}
	return nil
}

func ValidateWeekDateRange(initialDate, finalDate string) error {
	initial, err := parseWeekDate(initialDate)
	if err != nil {
		return ErrWeekInitialDateWrongFormat
	}
	final, err := parseWeekDate(finalDate)
	if err != nil {
		return ErrWeekFinalDateWrongFormat
	}
	if final.Sub(initial).Hours() != 24*6 {
		return ErrWeekDateRangeInvalid
	}
	return nil
}

func ValidateWeekCount(weeksCount int) error {
	switch weeksCount {
	case 8, 16:
		return nil
	default:
		return ErrWeekCountInvalid
	}
}

func isValidWeekDate(value string) bool {
	_, err := time.Parse(weekDateLayout, value)
	return err == nil
}

func parseWeekDate(value string) (time.Time, error) {
	return time.Parse(weekDateLayout, value)
}
