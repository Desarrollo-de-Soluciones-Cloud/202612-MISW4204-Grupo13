package domain_test

import (
	domain "backend/internal/weeks/domain"
	"errors"
	"testing"
)

func TestValidateWeekCount(t *testing.T) {
	if err := domain.ValidateWeekCount(8); err != nil {
		t.Fatalf("expected 8 weeks to be valid, got %v", err)
	}
	if err := domain.ValidateWeekCount(16); err != nil {
		t.Fatalf("expected 16 weeks to be valid, got %v", err)
	}
}

func TestValidateWeekCountRejectsInvalidValue(t *testing.T) {
	err := domain.ValidateWeekCount(12)
	if !errors.Is(err, domain.ErrWeekCountInvalid) {
		t.Fatalf("expected ErrWeekCountInvalid, got %v", err)
	}
}

func TestValidateWeekInitialDateIsMonday(t *testing.T) {
	if err := domain.ValidateWeekInitialDateIsMonday("2026-01-12"); err != nil {
		t.Fatalf("expected monday to be valid, got %v", err)
	}
}

func TestValidateWeekFinalDateIsSunday(t *testing.T) {
	if err := domain.ValidateWeekFinalDateIsSunday("2026-01-18"); err != nil {
		t.Fatalf("expected sunday to be valid, got %v", err)
	}
}
