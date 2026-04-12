package domain

import (
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

	switch {
	case trimmedName == "":
		return ErrPeriodNameRequired
	case len(trimmedName) != PeriodNameLength:
		return ErrPeriodNameWrongFormat
	default:
		return nil
	}
}

func ValidatePeriodInitialDate(date time.Time) error {
	if date.IsZero() {
		return ErrPeriodInitialDateRequired
	}
	return nil
}

func ValidatePeriodFinalDate(date time.Time) error {
	if date.IsZero() {
		return ErrPeriodFinalDateRequired
	}
	return nil
}

func ValidatePeriodInscriptionFinalDate(date time.Time) error {
	if date.IsZero() {
		return ErrPeriodInscriptionFinalDateRequired
	}
	return nil
}

func ValidatePeriodDateSequence(initialDate, finalDate, inscriptionFinalDate time.Time) error {
	if initialDate.After(finalDate) {
		return ErrPeriodDateSequenceInvalid
	}
	if initialDate.After(inscriptionFinalDate) {
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
