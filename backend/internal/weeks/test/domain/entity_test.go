package domain_test

import (
	domain "backend/internal/weeks/domain"
	"errors"
	"testing"
)

const testWeekInitialDate = "2026-01-12"

func TestNewWeekSuccess(t *testing.T) {
	week, err := domain.NewWeek(1, 1, testWeekInitialDate, "2026-01-18")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if week.Number != 1 {
		t.Fatalf("expected number 1, got %d", week.Number)
	}
	if week.PeriodID != 1 {
		t.Fatalf("expected period id 1, got %d", week.PeriodID)
	}
}

func TestNewWeekRejectsInvalidInitialDate(t *testing.T) {
	_, err := domain.NewWeek(1, 1, "2026-01-13", "2026-01-18")
	if !errors.Is(err, domain.ErrWeekInitialDateMustBeMonday) {
		t.Fatalf("expected ErrWeekInitialDateMustBeMonday, got %v", err)
	}
}

func TestNewWeekRejectsInvalidFinalDate(t *testing.T) {
	_, err := domain.NewWeek(1, 1, testWeekInitialDate, "2026-01-19")
	if !errors.Is(err, domain.ErrWeekFinalDateMustBeSunday) {
		t.Fatalf("expected ErrWeekFinalDateMustBeSunday, got %v", err)
	}
}

func TestNewWeekRejectsInvalidRange(t *testing.T) {
	_, err := domain.NewWeek(1, 1, testWeekInitialDate, "2026-01-25")
	if !errors.Is(err, domain.ErrWeekDateRangeInvalid) {
		t.Fatalf("expected ErrWeekDateRangeInvalid, got %v", err)
	}
}
