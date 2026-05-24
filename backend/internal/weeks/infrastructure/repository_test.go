package infrastructure

import (
	"backend/internal/shared/database"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupWeeksDryRunDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=test password=test dbname=test port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("expected dry run db, got %v", err)
	}

	database.DB = db
}

func TestNewWeekRepository(t *testing.T) {
	repo := NewWeekRepository()
	if repo == nil {
		t.Fatalf("expected week repository")
	}
}

func TestWeekRepositoryDryRunQueries(t *testing.T) {
	setupWeeksDryRunDB(t)

	repo := NewWeekRepository()
	if _, err := repo.FindAllByPeriodID(2); err != nil {
		t.Fatalf("expected no error listing weeks, got %v", err)
	}
	if _, err := repo.ExistsByPeriodID(2); err != nil {
		t.Fatalf("expected no error checking week existence, got %v", err)
	}
}
