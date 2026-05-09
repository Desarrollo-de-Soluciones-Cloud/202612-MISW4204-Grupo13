package application_test

import (
	application "backend/internal/weeks/application"
	"backend/internal/weeks/domain"
	"errors"
	"testing"
	"time"
)

func TestGetWeekStatusByPeriodAndNumberReturnsActive(t *testing.T) {
	repo := &mockWeekRepository{
		weeks: []domain.Week{
			{ID: 1, PeriodID: 1, Number: 2, InitialDate: "2026-01-19", FinalDate: "2026-01-25"},
		},
	}

	getWeekStatus := application.NewGetWeekStatusByPeriodAndNumberWithNow(repo, func() time.Time {
		return time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC)
	})

	output, err := getWeekStatus.Execute(application.GetWeekStatusByPeriodAndNumberInput{
		PeriodID: 1,
		Number:   2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Status != domain.WeekStatusActive {
		t.Fatalf("expected active status, got %s", output.Status)
	}
}

func TestGetWeekStatusByPeriodAndNumberReturnsClosed(t *testing.T) {
	repo := &mockWeekRepository{
		weeks: []domain.Week{
			{ID: 1, PeriodID: 1, Number: 1, InitialDate: "2026-01-12", FinalDate: "2026-01-18"},
		},
	}

	getWeekStatus := application.NewGetWeekStatusByPeriodAndNumberWithNow(repo, func() time.Time {
		return time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC)
	})

	output, err := getWeekStatus.Execute(application.GetWeekStatusByPeriodAndNumberInput{
		PeriodID: 1,
		Number:   1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Status != domain.WeekStatusClosed {
		t.Fatalf("expected closed status, got %s", output.Status)
	}
}

func TestGetWeekStatusByPeriodAndNumberRejectsInvalidNumber(t *testing.T) {
	repo := &mockWeekRepository{}
	getWeekStatus := application.NewGetWeekStatusByPeriodAndNumber(repo)

	_, err := getWeekStatus.Execute(application.GetWeekStatusByPeriodAndNumberInput{
		PeriodID: 1,
		Number:   0,
	})
	if !errors.Is(err, domain.ErrWeekNumberInvalid) {
		t.Fatalf("expected ErrWeekNumberInvalid, got %v", err)
	}
}
