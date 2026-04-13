package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"testing"
)

func TestListPeriodsSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate1 := "2024-01-01"
	finalDate1 := "2024-06-30"
	inscriptionDate1 := "2024-01-15"

    initialDate2 := "2024-07-01"
	finalDate2 := "2024-12-31"
	inscriptionDate2 := "2024-08-15"

	period1, _ := domain.NewPeriod("2024-01", initialDate1, finalDate1, inscriptionDate1, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2024-02", initialDate2, finalDate2, inscriptionDate2, domain.ClosedPeriod)

	mockRepo.Create(period1)
	mockRepo.Create(period2)

	listPeriods := applicationpkg.NewListPeriods(mockRepo)
	output, err := listPeriods.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Periods) != 2 {
		t.Errorf("expected 2 periods, got %d", len(output.Periods))
	}
}

func TestListPeriodsEmpty(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	listPeriods := applicationpkg.NewListPeriods(mockRepo)
	output, err := listPeriods.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Periods) != 0 {
		t.Errorf("expected 0 periods, got %d", len(output.Periods))
	}
}
