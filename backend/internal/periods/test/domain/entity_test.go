package domain

import (
	domainpkg "backend/internal/periods/domain"
	"errors"
	"strings"
	"testing"
)

const (
	testExpectedNoErrorMsg       = "expected no error, got %v"
	testExpectedNormalizedNameMsg = "expected normalized name, got %q"
)

func TestNewPeriodSuccess(t *testing.T) {
	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	period, err := domainpkg.NewPeriod(" "+testPeriodName202610+" ", initialDate, weeksCount, domainpkg.ActivePeriod)
	if err != nil {
		t.Fatalf(testExpectedNoErrorMsg, err)
	}

	if period.Name != testPeriodName202610 {
		t.Fatalf(testExpectedNormalizedNameMsg, period.Name)
	}
	if period.InitialDate != initialDate {
		t.Fatalf("expected initial date %v, got %v", initialDate, period.InitialDate)
	}
	if period.WeeksCount != weeksCount {
		t.Fatalf("expected weeks count %v, got %v", weeksCount, period.WeeksCount)
	}
	if period.PeriodState != domainpkg.ActivePeriod {
		t.Fatalf("expected state %q, got %q", domainpkg.ActivePeriod, period.PeriodState)
	}
}

func TestNewPeriodRejectsInvalidState(t *testing.T) {
	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	_, err := domainpkg.NewPeriod(testPeriodName202610, initialDate, weeksCount, domainpkg.PeriodState("invalid"))
	if !errors.Is(err, domainpkg.ErrPeriodStateInvalid) {
		t.Fatalf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestUpdatePeriodNormalizesValues(t *testing.T) {
	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	period, err := domainpkg.NewPeriod(testPeriodName202610, initialDate, weeksCount, domainpkg.ActivePeriod)
	if err != nil {
		t.Fatalf(testExpectedNoErrorMsg, err)
	}

	newInitialDate := testPeriodInitialDate1012
	newWeeksCount := 16

	err = period.UpdatePeriod(" 2026-11 ", newInitialDate, newWeeksCount, domainpkg.ClosedPeriod)
	if err != nil {
		t.Fatalf(testExpectedNoErrorMsg, err)
	}

	if period.Name != testPeriodName202611 {
		t.Fatalf(testExpectedNormalizedNameMsg, period.Name)
	}
	if period.PeriodState != domainpkg.ClosedPeriod {
		t.Fatalf("expected state %q, got %q", domainpkg.ClosedPeriod, period.PeriodState)
	}
}

func TestNewPeriodRejectsShortName(t *testing.T) {
	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	_, err := domainpkg.NewPeriod("", initialDate, weeksCount, domainpkg.ActivePeriod)
	if !errors.Is(err, domainpkg.ErrPeriodNameRequired) {
		t.Fatalf("expected ErrPeriodNameRequired, got %v", err)
	}
}

func TestNewPeriodRejectsNonMondayInitialDate(t *testing.T) {
	// 2026-10-06 is a Tuesday
	initialDate := "2026-10-06"
	weeksCount := 8

	_, err := domainpkg.NewPeriod(testPeriodName202610, initialDate, weeksCount, domainpkg.ActivePeriod)
	if !errors.Is(err, domainpkg.ErrPeriodInitialDateMustBeMonday) {
		t.Fatalf("expected ErrPeriodInitialDateMustBeMonday, got %v", err)
	}
}

func TestNewPeriodRejectsPastInitialDate(t *testing.T) {
	initialDate := "2020-11-16"
	weeksCount := 8

	_, err := domainpkg.NewPeriod("2020-10", initialDate, weeksCount, domainpkg.ActivePeriod)
	if !errors.Is(err, domainpkg.ErrPeriodInitialDateMustBeFuture) {
		t.Fatalf("expected ErrPeriodInitialDateMustBeFuture, got %v", err)
	}
}

func TestNewPeriodRejectsInvalidWeeksCount(t *testing.T) {
	initialDate := testPeriodInitialDate1005
	weeksCount := 10

	_, err := domainpkg.NewPeriod(testPeriodName202610, initialDate, weeksCount, domainpkg.ActivePeriod)
	if !errors.Is(err, domainpkg.ErrPeriodWeeksCountInvalid) {
		t.Fatalf("expected ErrPeriodWeeksCountInvalid, got %v", err)
	}
}

func TestUpdatePeriodRejectsInvalidState(t *testing.T) {
	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	period, err := domainpkg.NewPeriod(testPeriodName202610, initialDate, weeksCount, domainpkg.ActivePeriod)
	if err != nil {
		t.Fatalf(testExpectedNoErrorMsg, err)
	}

	err = period.UpdatePeriod(testPeriodName202610, initialDate, weeksCount, domainpkg.PeriodState("unknown"))
	if !errors.Is(err, domainpkg.ErrPeriodStateInvalid) {
		t.Fatalf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestNormalizePeriodName(t *testing.T) {
	normalized := domainpkg.NormalizePeriodName(" 2024-01 ")
	if normalized != "2024-01" {
		t.Fatalf(testExpectedNormalizedNameMsg, normalized)
	}
}

func TestValidatePeriodNameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 101)
	err := domainpkg.ValidatePeriodName(longName)
	if err != nil {
		t.Fatalf("expected no error for long name, got %v", err)
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
