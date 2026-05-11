package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"testing"
)

func TestListPeriodsSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate1 := testPeriodInitialDate1005
	weeksCount1 := 8

	initialDate2 := testPeriodInitialDate1012
	weeksCount2 := 16

	period1, _ := domain.NewPeriod(testPeriodName202610, initialDate1, weeksCount1, domain.ActivePeriod)
	period2, _ := domain.NewPeriod(testPeriodName202611, initialDate2, weeksCount2, domain.ClosedPeriod)

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

	if periodsByName[testPeriodName202610].WeeksCount != weeksCount1 {
		t.Errorf("expected period %s weeks count %d, got %d", testPeriodName202610, weeksCount1, periodsByName[testPeriodName202610].WeeksCount)
	}

	if periodsByName[testPeriodName202611].WeeksCount != weeksCount2 {
		t.Errorf("expected period %s weeks count %d, got %d", testPeriodName202611, weeksCount2, periodsByName[testPeriodName202611].WeeksCount)
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
