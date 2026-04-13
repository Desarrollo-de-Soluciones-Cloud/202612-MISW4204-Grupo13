package domain

import (
	"regexp"
	"strings"
	"time"
)

const (
	PeriodNameLength = 7
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

func ValidatePeriodFinalDate(date string) error {
	timmedDate := strings.TrimSpace(date)

	_, err := time.Parse("2006-01-02", date)
    if err != nil {
        return ErrPeriodFinalDateWrongFormat
    }

	switch {
	case timmedDate == "":
		return ErrPeriodFinalDateRequired
	case len(timmedDate) != 10:
		return ErrPeriodFinalDateWrongFormat
	}
	return nil
}

func ValidatePeriodInscriptionFinalDate(date string) error {
	timmedDate := strings.TrimSpace(date)

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrPeriodInscriptionFinalDateWrongFormat
	}

	switch {
	case timmedDate == "":
		return ErrPeriodInscriptionFinalDateRequired
	case len(timmedDate) != 10:
		return ErrPeriodInscriptionFinalDateWrongFormat
	}
	return nil
}

func ValidatePeriodDateSequence(initialDate, finalDate, inscriptionFinalDate string) error {
	if initialDate > finalDate {
		return ErrPeriodDateSequenceInvalid
	}
	if initialDate > inscriptionFinalDate {
		return ErrPeriodDateSequenceInvalid
	}
	if inscriptionFinalDate > finalDate {
		return ErrPeriodDateSequenceInvalid
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
