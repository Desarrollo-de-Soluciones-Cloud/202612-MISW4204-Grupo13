package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
	"time"
)

func TestListPeriodsByStateSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period1, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2024-02", initialDate.AddDate(0, 1, 0), finalDate.AddDate(0, 1, 0), inscriptionDate.AddDate(0, 1, 0), domain.ActivePeriod)
	period3, _ := domain.NewPeriod("2024-03", initialDate.AddDate(0, 2, 0), finalDate.AddDate(0, 2, 0), inscriptionDate.AddDate(0, 2, 0), domain.ClosedPeriod)

	mockRepo.Create(period1)
	mockRepo.Create(period2)
	mockRepo.Create(period3)

	listPeriodsByState := applicationpkg.NewListPeriodsByState(mockRepo)
	output, err := listPeriodsByState.Execute(applicationpkg.ListPeriodsByStateInput{
		PeriodState: domain.ActivePeriod,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Periods) != 2 {
		t.Errorf("expected 2 active periods, got %d", len(output.Periods))
	}
	for _, p := range output.Periods {
		if p.PeriodState != domain.ActivePeriod {
			t.Errorf("expected active period, got %q", p.PeriodState)
		}
	}
}

func TestListPeriodsByStateInvalidState(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	listPeriodsByState := applicationpkg.NewListPeriodsByState(mockRepo)

	_, err := listPeriodsByState.Execute(applicationpkg.ListPeriodsByStateInput{
		PeriodState: domain.PeriodState("invalid"),
	})

	if !errors.Is(err, domain.ErrPeriodStateInvalid) {
		t.Errorf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestListPeriodsByStateNoResults(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period1, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period1)

	listPeriodsByState := applicationpkg.NewListPeriodsByState(mockRepo)
	output, err := listPeriodsByState.Execute(applicationpkg.ListPeriodsByStateInput{
		PeriodState: domain.ClosedPeriod,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Periods) != 0 {
		t.Errorf("expected 0 closed periods, got %d", len(output.Periods))
	}
}
