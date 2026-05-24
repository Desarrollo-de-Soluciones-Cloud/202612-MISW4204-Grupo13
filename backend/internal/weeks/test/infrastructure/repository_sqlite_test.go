package infrastructure_test

import (
	sharedDB "backend/internal/shared/database/testsupport"
	backendweeks "backend/internal/weeks/domain"
	weeksinfra "backend/internal/weeks/infrastructure"
	"errors"
	"testing"
)

func TestWeekRepositorySQLiteCRUD(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &backendweeks.Week{})
	repo := weeksinfra.NewWeekRepository()

	weeks := []backendweeks.Week{
		{PeriodID: 1, Number: 1, InitialDate: "2026-04-06", FinalDate: "2026-04-12"},
		{PeriodID: 1, Number: 2, InitialDate: "2026-04-13", FinalDate: "2026-04-19"},
	}
	if err := repo.CreateMany(weeks); err != nil {
		t.Fatalf("expected create many, got %v", err)
	}

	items, err := repo.FindAllByPeriodID(1)
	if err != nil || len(items) != 2 {
		t.Fatalf("expected 2 weeks, got %v %d", err, len(items))
	}

	week, err := repo.FindByPeriodIDAndNumber(1, 1)
	if err != nil || week.Number != 1 {
		t.Fatalf("expected week 1, got %v %#v", err, week)
	}

	weekByDate, err := repo.FindByPeriodIDAndStartDate(1, "2026-04-06")
	if err != nil || weekByDate.Number != 1 {
		t.Fatalf("expected week by date, got %v %#v", err, weekByDate)
	}

	exists, err := repo.ExistsByPeriodID(1)
	if err != nil || !exists {
		t.Fatalf("expected exists, got %v %v", err, exists)
	}
}

func TestWeekRepositorySQLiteNotFound(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &backendweeks.Week{})
	repo := weeksinfra.NewWeekRepository()

	if _, err := repo.FindByID(999); !errors.Is(err, backendweeks.ErrWeekNotFound) {
		t.Fatalf("expected week not found, got %v", err)
	}
}
