package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestGetPeriodByNameSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"

	period, _ := domain.NewPeriod("2026-10", initialDate, 16, domain.ActivePeriod)
	mockRepo.Create(period)

	getPeriodByName := applicationpkg.NewGetPeriodByName(mockRepo)
	output, err := getPeriodByName.Execute(applicationpkg.GetPeriodByNameInput{Name: "2026-10"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2026-10" {
		t.Errorf("expected name '2026-10', got %q", output.Name)
	}
}

func TestGetPeriodByNameNotFound(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	getPeriodByName := applicationpkg.NewGetPeriodByName(mockRepo)

	_, err := getPeriodByName.Execute(applicationpkg.GetPeriodByNameInput{Name: "2024-10"})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}

func TestGetPeriodByNameInvalidName(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	getPeriodByName := applicationpkg.NewGetPeriodByName(mockRepo)

	_, err := getPeriodByName.Execute(applicationpkg.GetPeriodByNameInput{Name: "nonexistent"})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}
