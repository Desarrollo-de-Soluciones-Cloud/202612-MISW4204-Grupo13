package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestGetPeriodByNameSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := testPeriodInitialDate1005

	period, _ := domain.NewPeriod(testPeriodName202610, initialDate, 16, domain.ActivePeriod)
	mockRepo.Create(period)

	getPeriodByName := applicationpkg.NewGetPeriodByName(mockRepo)
	output, err := getPeriodByName.Execute(applicationpkg.GetPeriodByNameInput{Name: testPeriodName202610})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != testPeriodName202610 {
		t.Errorf("expected name %q, got %q", testPeriodName202610, output.Name)
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
