package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestGetPeriodByIDSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

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
