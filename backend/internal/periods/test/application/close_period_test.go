package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestClosePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	closePeriod := applicationpkg.NewClosePeriod(mockRepo)

	// Create an active period
	period, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	mockRepo.Create(period)

	input := applicationpkg.ClosePeriodInput{
		ID: period.ID,
	}

	output, err := closePeriod.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.PeriodState != domain.ClosedPeriod {
		t.Errorf("expected state %q, got %q", domain.ClosedPeriod, output.PeriodState)
	}

	if output.Name != "2026-10" {
		t.Errorf("expected name '2026-10', got %q", output.Name)
	}

	// Verify it was saved
	storedPeriod, _ := mockRepo.FindByID(output.ID)
	if storedPeriod.PeriodState != domain.ClosedPeriod {
		t.Errorf("expected stored period state %q, got %q", domain.ClosedPeriod, storedPeriod.PeriodState)
	}
}

func TestClosePeriodNotFound(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	closePeriod := applicationpkg.NewClosePeriod(mockRepo)

	input := applicationpkg.ClosePeriodInput{
		ID: 999,
	}

	_, err := closePeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}

func TestClosePeriodAlreadyClosed(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	closePeriod := applicationpkg.NewClosePeriod(mockRepo)

	// Create a closed period
	period, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ClosedPeriod)
	mockRepo.Create(period)

	input := applicationpkg.ClosePeriodInput{
		ID: period.ID,
	}

	_, err := closePeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodStateInvalid) {
		t.Errorf("expected ErrPeriodStateInvalid, got %v", err)
	}
}
