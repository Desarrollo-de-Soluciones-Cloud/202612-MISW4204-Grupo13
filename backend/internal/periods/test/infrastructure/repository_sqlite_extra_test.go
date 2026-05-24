package infrastructure_test

import (
	periodsDomain "backend/internal/periods/domain"
	periodsinfra "backend/internal/periods/infrastructure"
	sharedDB "backend/internal/shared/database/testsupport"
	"errors"
	"testing"
)

func TestPeriodRepositorySQLiteCRUD(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &periodsDomain.Period{})
	repo := periodsinfra.NewPeriodRepository()

	period := &periodsDomain.Period{
		Name: "2027-10", PeriodState: periodsDomain.ActivePeriod, InitialDate: "2027-01-05",
		FinalDate: "2027-05-02", WeeksCount: 16, InscriptionFinalDate: "2027-01-04",
	}
	if err := repo.Create(period); err != nil {
		t.Fatalf("expected create, got %v", err)
	}

	byID, err := repo.FindByID(period.ID)
	if err != nil || byID.Name != "2027-10" {
		t.Fatalf("expected find by id, got %v %#v", err, byID)
	}

	byName, err := repo.FindByName("2027-10")
	if err != nil || byName.ID != period.ID {
		t.Fatalf("expected find by name, got %v %#v", err, byName)
	}
}

func TestPeriodRepositorySQLiteDeleteNotFound(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &periodsDomain.Period{})
	repo := periodsinfra.NewPeriodRepository()

	if err := repo.Delete(999); !errors.Is(err, periodsDomain.ErrPeriodNotFound) {
		t.Fatalf("expected period not found, got %v", err)
	}
}
