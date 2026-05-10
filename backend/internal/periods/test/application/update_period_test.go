package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestUpdatePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	period, _ := domain.NewPeriod(testPeriodName202610, initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	output, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:   period.ID,
		Name: testPeriodName202611,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != testPeriodName202611 {
		t.Errorf("expected name %q, got %q", testPeriodName202611, output.Name)
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
		Name: testPeriodName202610,
	})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}

func TestUpdatePeriodInvalidName(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	period, _ := domain.NewPeriod(testPeriodName202610, initialDate, weeksCount, domain.ActivePeriod)
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

	initialDate := testPeriodInitialDate1005
	weeksCount := 8

	period1, _ := domain.NewPeriod(testPeriodName202610, initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period1)

	period2, _ := domain.NewPeriod(testPeriodName202611, initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period2)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:   period2.ID,
		Name: testPeriodName202610,
	})

	if !errors.Is(err, domain.ErrPeriodNameAlreadyExists) {
		t.Errorf("expected ErrPeriodNameAlreadyExists, got %v", err)
	}
}
