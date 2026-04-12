package infrastructure

import (
	"backend/internal/periods/domain"
	"backend/internal/periods/infrastructure"
	"backend/internal/shared/database"
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

func TestCreatePeriod(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	period := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-10-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	err := repo.Create(period)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if period.ID == 0 {
		t.Error("Create() did not assign an ID")
	}
}

func TestCreatePeriodDuplicateName(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	period1 := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-10-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	period2 := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-11-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	if err := repo.Create(period1); err != nil {
		t.Fatalf("Create() first period failed: %v", err)
	}

	err := repo.Create(period2)
	if err == nil {
		t.Error("Create() should fail with duplicate name constraint")
	}
}

func TestFindByID(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	period := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-10-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	if err := repo.Create(period); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	found, err := repo.FindByID(period.ID)
	if err != nil {
		t.Errorf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Error("FindByID() returned nil")
	}

	if found.Name != period.Name {
		t.Errorf("FindByID() name = %v, want %v", found.Name, period.Name)
	}
}

func TestFindByIDNotFound(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	_, err := repo.FindByID(999)
	if err != domain.ErrPeriodNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domain.ErrPeriodNotFound)
	}
}

func TestFindByName(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	period := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-10-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	if err := repo.Create(period); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	found, err := repo.FindByName("2026-10")
	if err != nil {
		t.Errorf("FindByName() error = %v", err)
	}

	if found == nil {
		t.Error("FindByName() returned nil")
	}

	if found.ID != period.ID {
		t.Errorf("FindByName() ID = %v, want %v", found.ID, period.ID)
	}
}

func TestFindByNameNotFound(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	_, err := repo.FindByName("nonexistent")
	if err != domain.ErrPeriodNotFound {
		t.Errorf("FindByName() error = %v, want %v", err, domain.ErrPeriodNotFound)
	}
}
func TestFindAll(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	periods := []domain.Period{
		{
			Name:                 "2026-10",
			InitialDate:          "2026-10-01",
			FinalDate:            "2026-12-31",
			InscriptionFinalDate: "2026-11-15",
			PeriodState:          domain.ActivePeriod,
		},
		{
			Name:                 "2026-01",
			InitialDate:          "2026-01-01",
			FinalDate:            "2026-04-30",
			InscriptionFinalDate: "2026-02-15",
			PeriodState:          domain.ActivePeriod,
		},
	}

	for i := range periods {
		if err := repo.Create(&periods[i]); err != nil {
			t.Fatalf("Create() failed: %v", err)
		}
	}

	found, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
	}

	if len(found) != 2 {
		t.Errorf("FindAll() returned %d periods, want 2", len(found))
	}
}

func TestFindAllEmpty(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	found, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll() error = %v", err)
	}

	if len(found) != 0 {
		t.Errorf("FindAll() returned %d periods, want 0", len(found))
	}
}

func TestFindAllByState(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	activePeriod := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-10-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	closedPeriod := &domain.Period{
		Name:                 "2026-01",
		InitialDate:          "2026-01-01",
		FinalDate:            "2026-04-30",
		InscriptionFinalDate: "2026-02-15",
		PeriodState:          domain.ClosedPeriod,
	}

	if err := repo.Create(activePeriod); err != nil {
		t.Fatalf("Create() active period failed: %v", err)
	}
	if err := repo.Create(closedPeriod); err != nil {
		t.Fatalf("Create() closed period failed: %v", err)
	}

	found, err := repo.FindAllByState(domain.ActivePeriod)
	if err != nil {
		t.Errorf("FindAllByState() error = %v", err)
	}

	if len(found) != 1 {
		t.Errorf("FindAllByState() returned %d periods, want 1", len(found))
	}

	if found[0].PeriodState != domain.ActivePeriod {
		t.Errorf("FindAllByState() state = %v, want %v", found[0].PeriodState, domain.ActivePeriod)
	}
}

func TestUpdate(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	period := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-10-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	if err := repo.Create(period); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	period.PeriodState = domain.ClosedPeriod
	err := repo.Update(period)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	found, err := repo.FindByID(period.ID)
	if err != nil {
		t.Fatalf("FindByID() failed: %v", err)
	}

	if found.PeriodState != domain.ClosedPeriod {
		t.Errorf("Update() state = %v, want %v", found.PeriodState, domain.ClosedPeriod)
	}
}

func TestDelete(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	period := &domain.Period{
		Name:                 "2026-10",
		InitialDate:          "2026-10-01",
		FinalDate:            "2026-12-31",
		InscriptionFinalDate: "2026-11-15",
		PeriodState:          domain.ActivePeriod,
	}

	if err := repo.Create(period); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err := repo.Delete(period.ID)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	_, err = repo.FindByID(period.ID)
	if err != domain.ErrPeriodNotFound {
		t.Errorf("Delete() did not delete the period: %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	setupTestDB(t)
	repo := infrastructure.NewPeriodRepository()

	err := repo.Delete(999)
	if err != domain.ErrPeriodNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domain.ErrPeriodNotFound)
	}
}
