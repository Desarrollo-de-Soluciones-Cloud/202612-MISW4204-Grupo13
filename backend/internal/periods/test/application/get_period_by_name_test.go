package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestGetPeriodByNameSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	getPeriodByName := applicationpkg.NewGetPeriodByName(mockRepo)
	output, err := getPeriodByName.Execute(applicationpkg.GetPeriodByNameInput{Name: "2024-01"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2024-01" {
		t.Errorf("expected name '2024-01', got %q", output.Name)
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

	_, err := getPeriodByName.Execute(applicationpkg.GetPeriodByNameInput{Name: "ab"})

	if !errors.Is(err, domain.ErrPeriodNameWrongFormat) {
		t.Errorf("expected ErrPeriodNameWrongFormat, got %v", err)
	}
}
