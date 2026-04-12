package domain

import (
	domainpkg "backend/internal/periods/domain"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewPeriodSuccess(t *testing.T) {
	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period, err := domainpkg.NewPeriod(" 2024-01 ", initialDate, finalDate, inscriptionDate, domainpkg.ActivePeriod)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if period.Name != "2024-01" {
		t.Fatalf("expected normalized name, got %q", period.Name)
	}
	if period.InitialDate != initialDate {
		t.Fatalf("expected initial date %v, got %v", initialDate, period.InitialDate)
	}
	if period.FinalDate != finalDate {
		t.Fatalf("expected final date %v, got %v", finalDate, period.FinalDate)
	}
	if period.InscriptionFinalDate != inscriptionDate {
		t.Fatalf("expected inscription date %v, got %v", inscriptionDate, period.InscriptionFinalDate)
	}
	if period.PeriodState != domainpkg.ActivePeriod {
		t.Fatalf("expected state %q, got %q", domainpkg.ActivePeriod, period.PeriodState)
	}
}

func TestNewPeriodRejectsInvalidState(t *testing.T) {
	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	_, err := domainpkg.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domainpkg.PeriodState("invalid"))
	if !errors.Is(err, domainpkg.ErrPeriodStateInvalid) {
		t.Fatalf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestUpdatePeriodNormalizesValues(t *testing.T) {
	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period, err := domainpkg.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domainpkg.ActivePeriod)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	newInitialDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	newFinalDate := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)
	newInscriptionDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	err = period.UpdatePeriod(" 2024-02 ", newInitialDate, newFinalDate, newInscriptionDate, domainpkg.ClosedPeriod)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if period.Name != "2024-02" {
		t.Fatalf("expected normalized name, got %q", period.Name)
	}
	if period.PeriodState != domainpkg.ClosedPeriod {
		t.Fatalf("expected state %q, got %q", domainpkg.ClosedPeriod, period.PeriodState)
	}
}

func TestNewPeriodRejectsShortName(t *testing.T) {
	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	_, err := domainpkg.NewPeriod("ab", initialDate, finalDate, inscriptionDate, domainpkg.ActivePeriod)
	if !errors.Is(err, domainpkg.ErrPeriodNameWrongFormat) {
		t.Fatalf("expected ErrPeriodNameWrongFormat, got %v", err)
	}
}

func TestNewPeriodRejectsInvalidDateSequence(t *testing.T) {
	initialDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	_, err := domainpkg.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domainpkg.ActivePeriod)
	if !errors.Is(err, domainpkg.ErrPeriodDateSequenceInvalid) {
		t.Fatalf("expected ErrPeriodDateSequenceInvalid, got %v", err)
	}
}

func TestNewPeriodRejectsInscriptionDateBeforeInitial(t *testing.T) {
	initialDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := domainpkg.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domainpkg.ActivePeriod)
	if !errors.Is(err, domainpkg.ErrPeriodDateSequenceInvalid) {
		t.Fatalf("expected ErrPeriodDateSequenceInvalid, got %v", err)
	}
}

func TestUpdatePeriodRejectsInvalidState(t *testing.T) {
	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period, err := domainpkg.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domainpkg.ActivePeriod)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = period.UpdatePeriod("2024-01", initialDate, finalDate, inscriptionDate, domainpkg.PeriodState("unknown"))
	if !errors.Is(err, domainpkg.ErrPeriodStateInvalid) {
		t.Fatalf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestNormalizePeriodName(t *testing.T) {
	normalized := domainpkg.NormalizePeriodName(" 2024-01 ")
	if normalized != "2024-01" {
		t.Fatalf("expected normalized name, got %q", normalized)
	}
}

func TestValidatePeriodNameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 101)
	err := domainpkg.ValidatePeriodName(longName)
	if !errors.Is(err, domainpkg.ErrPeriodNameWrongFormat) {
		t.Fatalf("expected error for too long name, got %v", err)
	}
}

func TestValidatePeriodDateSequence(t *testing.T) {
	initialDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	err := domainpkg.ValidatePeriodDateSequence(initialDate, finalDate, inscriptionDate)
	if !errors.Is(err, domainpkg.ErrPeriodDateSequenceInvalid) {
		t.Fatalf("expected ErrPeriodDateSequenceInvalid, got %v", err)
	}
}

func TestIsValidPeriodState(t *testing.T) {
	if !domainpkg.IsValidPeriodState(domainpkg.ActivePeriod) {
		t.Fatalf("expected ActivePeriod to be valid")
	}
	if !domainpkg.IsValidPeriodState(domainpkg.ClosedPeriod) {
		t.Fatalf("expected ClosedPeriod to be valid")
	}
	if domainpkg.IsValidPeriodState(domainpkg.PeriodState("invalid")) {
		t.Fatalf("expected invalid state to be rejected")
	}
}
