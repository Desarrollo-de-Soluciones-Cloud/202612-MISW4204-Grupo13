package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestGetPeriodByIDSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	period, _ := domain.NewPeriod(testPeriodName202610, initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	getPeriodByID := applicationpkg.NewGetPeriodByID(mockRepo)
	output, err := getPeriodByID.Execute(applicationpkg.GetPeriodByIDInput{ID: period.ID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != testPeriodName202610 {
		t.Errorf("expected name %q, got %q", testPeriodName202610, output.Name)
	}
	if output.PeriodState != domain.ActivePeriod {
		t.Errorf("expected state %q, got %q", domain.ActivePeriod, output.PeriodState)
	}
	if output.WeeksCount != weeksCount {
		t.Errorf("expected weeks count %d, got %d", weeksCount, output.WeeksCount)
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
