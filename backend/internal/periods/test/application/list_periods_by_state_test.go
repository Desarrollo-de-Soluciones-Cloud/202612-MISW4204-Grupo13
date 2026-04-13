package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestListPeriodsByStateSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate1 := "2026-10-05"
	initialDate2 := "2026-10-12"
	initialDate3 := "2026-11-09"

	period1, _ := domain.NewPeriod("2026-10", initialDate1, 16, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2026-11", initialDate2, 16, domain.ActivePeriod)
	period3, _ := domain.NewPeriod("2026-12", initialDate3, 16, domain.ClosedPeriod)

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

		initialDate := "2026-10-05"

		period1, _ := domain.NewPeriod("2026-10", initialDate, 16, domain.ActivePeriod)
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
