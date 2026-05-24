package infrastructure_test

import (
	"testing"

	"backend/internal/periods/domain"
	infrastructurepkg "backend/internal/periods/infrastructure"
	"backend/internal/shared/database"
)

const testFindAllByStateErrMsg = "FindAllByState() error = %v, want nil"

// setupTestDB initializes test database (skips if not available)
func setupTestDB(t *testing.T) {
	// Skip infrastructure tests if database not initialized
	// Database operations are tested at application layer with mocks
	// This file documents infrastructure repository interface
	if database.DB == nil {
		t.Skip("database.DB not initialized - infrastructure tests require database setup")
	}
	
	// Clean up before each test
	if err := database.DB.Exec("DELETE FROM periods").Error; err != nil {
		t.Logf("warning: could not clean periods table: %v", err)
	}
}

// TestCreatePeriod tests the Create method of PeriodRepository
func TestCreate(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)

	err := repo.Create(period)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	if period.ID == 0 {
		t.Errorf("Create() did not assign ID to period")
	}
}

// TestFindByID tests the FindByID method
func TestFindByID(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	found, err := repo.FindByID(period.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v, want nil", err)
	}

	if found.Name != testPeriodName202610 {
		t.Errorf("FindByID() name = %s, want %s", found.Name, testPeriodName202610)
	}
	if found.WeeksCount != 8 {
		t.Errorf("FindByID() weeks count = %d, want 8", found.WeeksCount)
	}
}

// TestFindByIDNotFound tests when period doesn't exist
func TestFindByIDNotFound(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	_, err := repo.FindByID(999)
	if err != domain.ErrPeriodNotFound {
		t.Errorf("FindByID() error = %v, want %v", err, domain.ErrPeriodNotFound)
	}
}

// TestFindByName tests the FindByName method
func TestFindByName(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	found, err := repo.FindByName(testPeriodName202610)
	if err != nil {
		t.Fatalf("FindByName() error = %v, want nil", err)
	}

	if found.Name != testPeriodName202610 {
		t.Errorf("FindByName() name = %s, want %s", found.Name, testPeriodName202610)
	}
}

// TestFindByNameNotFound tests when period with that name doesn't exist
func TestFindByNameNotFound(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	_, err := repo.FindByName("nonexistent")
	if err != domain.ErrPeriodNotFound {
		t.Errorf("FindByName() error = %v, want %v", err, domain.ErrPeriodNotFound)
	}
}

// TestFindAll tests the FindAll method
func TestFindAll(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period1, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	period2, _ := domain.NewPeriod(testPeriodName202611, testPeriodInitialDate1012, 16, domain.ClosedPeriod)
	repo.Create(period1)
	repo.Create(period2)

	periods, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error = %v, want nil", err)
	}

	if len(periods) != 2 {
		t.Errorf("FindAll() returned %d periods, want 2", len(periods))
	}
}

// TestFindAllEmpty tests FindAll with no periods
func TestFindAllEmpty(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	periods, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error = %v, want nil", err)
	}

	if len(periods) != 0 {
		t.Errorf("FindAll() returned %d periods, want 0", len(periods))
	}
}

// TestFindAllByState tests fetching periods by state
func TestFindAllByState(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period1, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	period2, _ := domain.NewPeriod(testPeriodName202611, testPeriodInitialDate1012, 16, domain.ActivePeriod)
	period3, _ := domain.NewPeriod("2026-12", "2026-11-09", 8, domain.ClosedPeriod)
	repo.Create(period1)
	repo.Create(period2)
	repo.Create(period3)

	activePeriods, err := repo.FindAllByState(domain.ActivePeriod)
	if err != nil {
		t.Fatalf(testFindAllByStateErrMsg, err)
	}

	if len(activePeriods) != 2 {
		t.Errorf("FindAllByState() returned %d active periods, want 2", len(activePeriods))
	}

	closedPeriods, err := repo.FindAllByState(domain.ClosedPeriod)
	if err != nil {
		t.Fatalf(testFindAllByStateErrMsg, err)
	}

	if len(closedPeriods) != 1 {
		t.Errorf("FindAllByState() returned %d closed periods, want 1", len(closedPeriods))
	}
}

// TestFindAllByStateEmpty tests FindAllByState with no results
func TestFindAllByStateEmpty(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period1, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period1)

	closedPeriods, err := repo.FindAllByState(domain.ClosedPeriod)
	if err != nil {
		t.Fatalf(testFindAllByStateErrMsg, err)
	}

	if len(closedPeriods) != 0 {
		t.Errorf("FindAllByState() returned %d closed periods, want 0", len(closedPeriods))
	}
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	err := period.UpdatePeriod(testPeriodName202610, testPeriodInitialDate1012, 16, domain.ClosedPeriod)
	if err != nil {
		t.Fatalf("UpdatePeriod() error = %v, want nil", err)
	}

	err = repo.Update(period)
	if err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}

	updated, _ := repo.FindByID(period.ID)
	if updated.WeeksCount != 16 {
		t.Errorf("Update() weeks count = %d, want 16", updated.WeeksCount)
	}
	if updated.PeriodState != domain.ClosedPeriod {
		t.Errorf("Update() state = %s, want %s", updated.PeriodState, domain.ClosedPeriod)
	}
}

// TestUpdateNotFound tests updating a non-existent period
func TestUpdateNotFound(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	fakePeriod := &domain.Period{
		ID:                   999,
		Name:                 testPeriodName202610,
		InitialDate:          testPeriodInitialDate1005,
		FinalDate:            "2026-11-29",
		InscriptionFinalDate: "2026-10-04",
		WeeksCount:           8,
		PeriodState:          domain.ActivePeriod,
	}

	err := repo.Update(fakePeriod)
	if err != nil {
		t.Errorf("Update() error = %v, want nil", err)
	}
}

// TestDelete tests the Delete method
func TestDelete(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	period, _ := domain.NewPeriod(testPeriodName202610, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	err := repo.Delete(period.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	_, err = repo.FindByID(period.ID)
	if err != domain.ErrPeriodNotFound {
		t.Errorf("Delete() did not remove the period")
	}
}

// TestDeleteNotFound tests deleting a non-existent period
func TestDeleteNotFound(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	err := repo.Delete(999)
	if err != domain.ErrPeriodNotFound {
		t.Errorf("Delete() error = %v, want %v", err, domain.ErrPeriodNotFound)
	}
}

// TestAutoMigrate tests the AutoMigrate method
func TestAutoMigrate(t *testing.T) {
	setupTestDB(t)
	repo := infrastructurepkg.NewPeriodRepository()

	err := repo.AutoMigrate()
	if err != nil {
		t.Fatalf("AutoMigrate() error = %v, want nil", err)
	}

	// Verify the table exists by checking schema
	if !database.DB.Migrator().HasTable(&domain.Period{}) {
		t.Errorf("AutoMigrate() did not create periods table")
	}
}

// TestNewPeriodRepository tests the factory function
func TestNewPeriodRepository(t *testing.T) {
	repo := infrastructurepkg.NewPeriodRepository()
	if repo == nil {
		t.Fatalf("NewPeriodRepository() returned nil")
	}
}
