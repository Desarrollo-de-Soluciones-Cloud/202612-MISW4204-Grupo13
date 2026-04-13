package infrastructure_test

import (
	"backend/internal/shared/database"
	weeksDomain "backend/internal/weeks/domain"
	weeksInfrastructure "backend/internal/weeks/infrastructure"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db
}

func TestWeekRepositoryCreateManyAndFindByPeriodID(t *testing.T) {
	setupTestDB(t)

	repo := weeksInfrastructure.NewWeekRepository()
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	weeks := []weeksDomain.Week{
		{PeriodID: 1, Number: 1, InitialDate: "2026-01-12", FinalDate: "2026-01-18"},
		{PeriodID: 1, Number: 2, InitialDate: "2026-01-19", FinalDate: "2026-01-25"},
	}

	if err := repo.CreateMany(weeks); err != nil {
		t.Fatalf("expected create many, got %v", err)
	}

	exists, err := repo.ExistsByPeriodID(1)
	if err != nil {
		t.Fatalf("expected exists query, got %v", err)
	}
	if !exists {
		t.Fatalf("expected period 1 to have weeks")
	}

	foundWeeks, err := repo.FindAllByPeriodID(1)
	if err != nil {
		t.Fatalf("expected find by period id, got %v", err)
	}
	if len(foundWeeks) != 2 {
		t.Fatalf("expected 2 weeks, got %d", len(foundWeeks))
	}
	if foundWeeks[0].Number != 1 || foundWeeks[1].Number != 2 {
		t.Fatalf("expected ordered weeks, got %+v", foundWeeks)
	}
}
