package infrastructure

import (
	"backend/internal/shared/database"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupReportsDryRunDB(t *testing.T) {
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

func TestNewAssignmentReader(t *testing.T) {
	reader := NewAssignmentReader()
	if reader == nil {
		t.Fatalf("expected assignment reader")
	}
}

func TestAssignmentReaderFindAllByWorkspaceID(t *testing.T) {
	setupReportsDryRunDB(t)

	reader := NewAssignmentReader()
	assignments, err := reader.FindAllByWorkspaceID(10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if assignments == nil {
		t.Fatalf("expected assignments slice")
	}
}
