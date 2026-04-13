package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestListPeriodsByStateSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate1 := "2024-01-01"
	finalDate1 := "2024-06-30"
	inscriptionDate1 := "2024-01-15"

    initialDate2 := "2024-07-01"
	finalDate2 := "2024-12-31"
	inscriptionDate2 := "2024-08-15"

	initialDate3 := "2025-01-01"
	finalDate3 := "2025-06-30"
	inscriptionDate3 := "2025-01-15"

	period1, _ := domain.NewPeriod("2024-01", initialDate1, finalDate1, inscriptionDate1, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2024-02", initialDate2, finalDate2, inscriptionDate2, domain.ActivePeriod)
	period3, _ := domain.NewPeriod("2024-03", initialDate3, finalDate3, inscriptionDate3, domain.ClosedPeriod)

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

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

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
