package domain

import (
	"regexp"
	"strings"
	"time"
)

const (
	PeriodNameLength = 7
	WeeksCount8      = 8
	WeeksCount16     = 16
)

func NormalizePeriodName(name string) string {
	return strings.TrimSpace(name)
}

func ValidatePeriodName(name string) error {
	trimmedName := strings.TrimSpace(name)

	matched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, trimmedName)
	if !matched {
		return ErrPeriodNameWrongFormat
	}

	switch {
	case trimmedName == "":
		return ErrPeriodNameRequired
	case len(trimmedName) != PeriodNameLength:
		return ErrPeriodNameWrongFormat
	default:
		return nil
	}
}

func ValidatePeriodInitialDate(date string) error {
	timmedDate := strings.TrimSpace(date)

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrPeriodInitialDateWrongFormat
	}

	switch {
	case timmedDate == "":
		return ErrPeriodInitialDateRequired
	case len(timmedDate) != 10:
		return ErrPeriodInitialDateWrongFormat
	}
	return nil
}

func ValidatePeriodInitialDateIsMonday(date string) error {
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	if err != nil {
		return ErrPeriodInitialDateWrongFormat
	}

	if parsedDate.Weekday() != time.Monday {
		return ErrPeriodInitialDateMustBeMonday
	}

	return nil
}

func ValidatePeriodInitialDateIsFuture(date string) error {
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	if err != nil {
		return ErrPeriodInitialDateWrongFormat
	}

	now := time.Now()
	if parsedDate.Before(now) || parsedDate.Format("2006-01-02") == now.Format("2006-01-02") {
		return ErrPeriodInitialDateMustBeFuture
	}

	return nil
}

func ValidatePeriodWeeksCount(weeksCount int) error {
	if weeksCount != WeeksCount8 && weeksCount != WeeksCount16 {
		return ErrPeriodWeeksCountInvalid
	}
	return nil
}

func ValidatePeriodState(state PeriodState) error {
	switch {
	case strings.TrimSpace(string(state)) == "":
		return ErrPeriodStateRequired
	case !IsValidPeriodState(state):
		return ErrPeriodStateInvalid
	default:
		return nil
	}
}

// Helper functions for date calculations

func CalculatePeriodFinalDate(initialDate string, weeksCount int) (string, error) {
	// Validate initial date is Monday
	if err := ValidatePeriodInitialDateIsMonday(initialDate); err != nil {
		return "", err
	}

	// Validate weeks count
	if err := ValidatePeriodWeeksCount(weeksCount); err != nil {
		return "", err
	}

	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(initialDate))
	if err != nil {
		return "", ErrPeriodInitialDateWrongFormat
	}

	// final_date = initial_date + (weeks_count * 7) - 1 day (to make it Sunday)
	finalDate := parsedDate.AddDate(0, 0, (weeksCount*7)-1)
	return finalDate.Format("2006-01-02"), nil
}

func CalculatePeriodInscriptionFinalDate(initialDate string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(initialDate))
	if err != nil {
		return "", ErrPeriodInitialDateWrongFormat
	}

	// inscription_final_date = initial_date - 1 day
	inscriptionFinalDate := parsedDate.AddDate(0, 0, -1)
	return inscriptionFinalDate.Format("2006-01-02"), nil
}

