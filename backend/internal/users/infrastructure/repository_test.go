package infrastructure

import (
	"backend/internal/shared/database"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUsersDryRunDB(t *testing.T) {
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

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository()
	if repo == nil {
		t.Fatalf("expected user repository")
	}
}

func TestUserRepositoryDryRunQueries(t *testing.T) {
	setupUsersDryRunDB(t)

	repo := NewUserRepository()
	if _, err := repo.FindAll(); err != nil {
		t.Fatalf("expected no error listing users, got %v", err)
	}
	if _, err := repo.FindAllByRole("admin"); err != nil {
		t.Fatalf("expected no error listing by role, got %v", err)
	}
}
