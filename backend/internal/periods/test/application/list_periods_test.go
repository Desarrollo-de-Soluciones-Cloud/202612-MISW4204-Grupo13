package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"testing"
	"time"
)

func TestListPeriodsSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period1, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2024-02", initialDate.AddDate(0, 1, 0), finalDate.AddDate(0, 1, 0), inscriptionDate.AddDate(0, 1, 0), domain.ClosedPeriod)

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
