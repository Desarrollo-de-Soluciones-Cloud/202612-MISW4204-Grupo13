package application_test

import (
	application "backend/internal/weeks/application"
	"backend/internal/weeks/domain"
	"errors"
	"testing"
)

func TestGetWeekByPeriodAndNumberSuccess(t *testing.T) {
	repo := &MockWeekRepository{
		weeks: []domain.Week{
			{ID: 1, PeriodID: 1, Number: 1, InitialDate: "2026-01-12", FinalDate: "2026-01-18"},
		},
	}

	getWeek := application.NewGetWeekByPeriodAndNumber(repo)
	output, err := getWeek.Execute(application.GetWeekByPeriodAndNumberInput{
		PeriodID: 1,
		Number:   1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Number != 1 || output.PeriodID != 1 {
		t.Fatalf("unexpected output: %+v", output)
	}
}

func TestGetWeekByPeriodAndNumberNotFound(t *testing.T) {
	repo := &MockWeekRepository{}
	getWeek := application.NewGetWeekByPeriodAndNumber(repo)

	_, err := getWeek.Execute(application.GetWeekByPeriodAndNumberInput{
		PeriodID: 1,
		Number:   1,
	})
	if !errors.Is(err, domain.ErrWeekNotFound) {
		t.Fatalf("expected ErrWeekNotFound, got %v", err)
	}
}

func TestGetWeekByPeriodAndNumberRejectsInvalidNumber(t *testing.T) {
	repo := &MockWeekRepository{}
	getWeek := application.NewGetWeekByPeriodAndNumber(repo)

	_, err := getWeek.Execute(application.GetWeekByPeriodAndNumberInput{
		PeriodID: 1,
		Number:   0,
	})
	if !errors.Is(err, domain.ErrWeekNumberInvalid) {
		t.Fatalf("expected ErrWeekNumberInvalid, got %v", err)
	}
}
