package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
	"time"
)

func TestGetPeriodByNameSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

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

	_, err := getPeriodByName.Execute(applicationpkg.GetPeriodByNameInput{Name: "2024-99"})

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
