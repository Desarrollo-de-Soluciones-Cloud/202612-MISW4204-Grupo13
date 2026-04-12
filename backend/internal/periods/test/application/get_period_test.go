package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
	"time"
)

func TestGetPeriodByIDSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	getPeriodByID := applicationpkg.NewGetPeriodByID(mockRepo)
	output, err := getPeriodByID.Execute(applicationpkg.GetPeriodByIDInput{ID: period.ID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2024-01" {
		t.Errorf("expected name '2024-01', got %q", output.Name)
	}
	if output.PeriodState != domain.ActivePeriod {
		t.Errorf("expected state %q, got %q", domain.ActivePeriod, output.PeriodState)
	}
}

func TestGetPeriodByIDNotFound(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	getPeriodByID := applicationpkg.NewGetPeriodByID(mockRepo)

	_, err := getPeriodByID.Execute(applicationpkg.GetPeriodByIDInput{ID: 999})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}
