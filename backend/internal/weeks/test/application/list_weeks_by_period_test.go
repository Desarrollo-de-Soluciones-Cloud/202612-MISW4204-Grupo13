package application_test

import (
	application "backend/internal/weeks/application"
	"backend/internal/weeks/domain"
	"errors"
	"testing"
)

func TestListWeeksByPeriodSuccess(t *testing.T) {
	repo := &MockWeekRepository{
		weeks: []domain.Week{
			{ID: 1, PeriodID: 1, Number: 1, InitialDate: "2026-01-12", FinalDate: "2026-01-18"},
			{ID: 2, PeriodID: 1, Number: 2, InitialDate: "2026-01-19", FinalDate: "2026-01-25"},
			{ID: 3, PeriodID: 2, Number: 1, InitialDate: "2026-02-02", FinalDate: "2026-02-08"},
		},
	}

	listWeeks := application.NewListWeeksByPeriod(repo)
	output, err := listWeeks.Execute(application.ListWeeksByPeriodInput{PeriodID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Weeks) != 2 {
		t.Fatalf("expected 2 weeks, got %d", len(output.Weeks))
	}
}

func TestListWeeksByPeriodRejectsInvalidPeriodID(t *testing.T) {
	repo := &MockWeekRepository{}
	listWeeks := application.NewListWeeksByPeriod(repo)

	_, err := listWeeks.Execute(application.ListWeeksByPeriodInput{PeriodID: 0})
	if !errors.Is(err, domain.ErrWeekPeriodIDRequired) {
		t.Fatalf("expected ErrWeekPeriodIDRequired, got %v", err)
	}
}
