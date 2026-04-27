package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestUpdatePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	output, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:   period.ID,
		Name: "2026-11",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2026-11" {
		t.Errorf("expected name '2026-11', got %q", output.Name)
	}
	if output.PeriodState != domain.ActivePeriod {
		t.Errorf("expected state to remain %q, got %q", domain.ActivePeriod, output.PeriodState)
	}
	if output.WeeksCount != weeksCount {
		t.Errorf("expected weeks count to remain %d, got %d", weeksCount, output.WeeksCount)
	}
}

func TestUpdatePeriodNotFound(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:   999,
		Name: "2026-10",
	})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}

func TestUpdatePeriodInvalidName(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:   period.ID,
		Name: "",
	})

	if !errors.Is(err, domain.ErrPeriodNameRequired) {
		t.Errorf("expected ErrPeriodNameRequired, got %v", err)
	}
}

func TestUpdatePeriodNameAlreadyExists(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period1, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period1)

	period2, _ := domain.NewPeriod("2026-11", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period2)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:   period2.ID,
		Name: "2026-10",
	})

	if !errors.Is(err, domain.ErrPeriodNameAlreadyExists) {
		t.Errorf("expected ErrPeriodNameAlreadyExists, got %v", err)
	}
}
