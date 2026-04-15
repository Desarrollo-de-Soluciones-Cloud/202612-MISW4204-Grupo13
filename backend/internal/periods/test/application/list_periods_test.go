package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"testing"
)

func TestListPeriodsSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate1 := "2026-10-05"
	weeksCount1 := 8

	initialDate2 := "2026-10-12"
	weeksCount2 := 16

	period1, _ := domain.NewPeriod("2026-10", initialDate1, weeksCount1, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2026-11", initialDate2, weeksCount2, domain.ClosedPeriod)

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

	periodsByName := make(map[string]applicationpkg.PeriodDTO, len(output.Periods))
	for _, period := range output.Periods {
		periodsByName[period.Name] = period
	}

	if periodsByName["2026-10"].WeeksCount != weeksCount1 {
		t.Errorf("expected period 2026-10 weeks count %d, got %d", weeksCount1, periodsByName["2026-10"].WeeksCount)
	}

	if periodsByName["2026-11"].WeeksCount != weeksCount2 {
		t.Errorf("expected period 2026-11 weeks count %d, got %d", weeksCount2, periodsByName["2026-11"].WeeksCount)
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
