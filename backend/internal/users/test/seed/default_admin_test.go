package seed

import (
	"backend/internal/shared/database"
	usersInfrastructure "backend/internal/users/infrastructure"
	usersSeed "backend/internal/users/seed"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSeedTestDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db
}

func TestSeedDefaultAdminCreatesUserAndIsIdempotent(t *testing.T) {
	setupSeedTestDB(t)

	t.Setenv("DEFAULT_ADMIN_NAME", "System Administrator")
	t.Setenv("DEFAULT_ADMIN_EMAIL", "admin@seneprojects.local")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "SeneprojectsAdmin2026!")

	if err := usersSeed.SeedDefaultAdmin(); err != nil {
		t.Fatalf("expected seed admin, got %v", err)
	}
	if err := usersSeed.SeedDefaultAdmin(); err != nil {
		t.Fatalf("expected seed admin to be idempotent, got %v", err)
	}

	repo := usersInfrastructure.NewUserRepository()
	user, err := repo.FindByEmail("admin@seneprojects.local")
	if err != nil {
		t.Fatalf("expected seeded admin, got %v", err)
	}
	if user.Name != "System Administrator" {
		t.Fatalf("expected seeded name, got %q", user.Name)
	}
}

func TestSeedDefaultAdminRequiresEnvValues(t *testing.T) {
	setupSeedTestDB(t)

	t.Setenv("DEFAULT_ADMIN_NAME", "")
	t.Setenv("DEFAULT_ADMIN_EMAIL", "admin@seneprojects.local")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "SeneprojectsAdmin2026!")

	if err := usersSeed.SeedDefaultAdmin(); err == nil {
		t.Fatal("expected missing env error")
	}
}

func TestSeedDefaultAdminRequiresEmail(t *testing.T) {
	setupSeedTestDB(t)

	t.Setenv("DEFAULT_ADMIN_NAME", "System Administrator")
	t.Setenv("DEFAULT_ADMIN_EMAIL", "")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "SeneprojectsAdmin2026!")

	if err := usersSeed.SeedDefaultAdmin(); err == nil {
		t.Fatal("expected missing email env error")
	}
}

func TestSeedDefaultAdminRequiresPassword(t *testing.T) {
	setupSeedTestDB(t)

	t.Setenv("DEFAULT_ADMIN_NAME", "System Administrator")
	t.Setenv("DEFAULT_ADMIN_EMAIL", "admin@seneprojects.local")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "")

	if err := usersSeed.SeedDefaultAdmin(); err == nil {
		t.Fatal("expected missing password env error")
	}
}

func TestSeedDefaultAdminPropagatesCreateValidationError(t *testing.T) {
	setupSeedTestDB(t)

	t.Setenv("DEFAULT_ADMIN_NAME", "System Administrator")
	t.Setenv("DEFAULT_ADMIN_EMAIL", "admin@seneprojects.local")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "short")

	if err := usersSeed.SeedDefaultAdmin(); err == nil {
		t.Fatal("expected create validation error")
	}
}

func TestSeedDefaultAdminPropagatesAutoMigrateError(t *testing.T) {
	setupSeedTestDB(t)

	sqlDB, err := database.DB.DB()
	if err != nil {
		t.Fatalf("expected sql db, got %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("expected closed db, got %v", err)
	}

	t.Setenv("DEFAULT_ADMIN_NAME", "System Administrator")
	t.Setenv("DEFAULT_ADMIN_EMAIL", "admin@seneprojects.local")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "SeneprojectsAdmin2026!")

	if err := usersSeed.SeedDefaultAdmin(); err == nil {
		t.Fatal("expected automigrate error")
	}
}
