package seed

import (
	"backend/internal/shared/database"
	usersInfrastructure "backend/internal/users/infrastructure"
	usersSeed "backend/internal/users/seed"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSeedDefaultAdminCreatesUserAndIsIdempotent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

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
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("expected sqlite db, got %v", err)
	}
	database.DB = db

	t.Setenv("DEFAULT_ADMIN_NAME", "")
	t.Setenv("DEFAULT_ADMIN_EMAIL", "admin@seneprojects.local")
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "SeneprojectsAdmin2026!")

	if err := usersSeed.SeedDefaultAdmin(); err == nil {
		t.Fatal("expected missing env error")
	}
}
