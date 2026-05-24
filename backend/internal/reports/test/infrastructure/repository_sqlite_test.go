package infrastructure_test

import (
	reportsDomain "backend/internal/reports/domain"
	reportsinfra "backend/internal/reports/infrastructure"
	sharedDB "backend/internal/shared/database/testsupport"
	"errors"
	"testing"
)

func TestReportRepositorySQLiteCRUD(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &reportsDomain.Report{})
	repo := reportsinfra.NewReportRepository()

	report, err := reportsDomain.NewWeeklyReport(1, 2, 3, 4, "reports/file.pdf")
	if err != nil {
		t.Fatalf("expected report, got %v", err)
	}
	if err := repo.Save(report); err != nil {
		t.Fatalf("expected save, got %v", err)
	}

	found, err := repo.FindByID(report.ID)
	if err != nil || found.FilePath != "reports/file.pdf" {
		t.Fatalf("expected find by id, got %v %#v", err, found)
	}
}

func TestReportRepositorySQLiteNotFound(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &reportsDomain.Report{})
	repo := reportsinfra.NewReportRepository()

	if _, err := repo.FindByID(999); !errors.Is(err, reportsDomain.ErrReportNotFound) {
		t.Fatalf("expected report not found, got %v", err)
	}
}
