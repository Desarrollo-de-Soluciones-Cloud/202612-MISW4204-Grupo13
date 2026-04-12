package infrastructure

import (
	"backend/internal/shared/database"
	usersDomain "backend/internal/users/domain"
	usersInfrastructure "backend/internal/users/infrastructure"
	"errors"
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

func TestUserRepositoryCRUD(t *testing.T) {
	setupTestDB(t)

	repo := usersInfrastructure.NewUserRepository()
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	user, err := usersDomain.NewUser("John Doe", "john@example.com", "password123", usersDomain.RoleProfessor)
	if err != nil {
		t.Fatalf("expected new user, got %v", err)
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("expected create user, got %v", err)
	}

	foundByID, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("expected find by id, got %v", err)
	}
	if foundByID.Email != "john@example.com" {
		t.Fatalf("expected email john@example.com, got %q", foundByID.Email)
	}

	foundByEmail, err := repo.FindByEmail("john@example.com")
	if err != nil {
		t.Fatalf("expected find by email, got %v", err)
	}
	if foundByEmail.ID != user.ID {
		t.Fatalf("expected same user id, got %d", foundByEmail.ID)
	}

	foundByID.Name = "John Updated"
	if err := repo.Update(foundByID); err != nil {
		t.Fatalf("expected update user, got %v", err)
	}

	users, err := repo.FindAll()
	if err != nil {
		t.Fatalf("expected find all, got %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	usersByRole, err := repo.FindAllByRole(usersDomain.RoleProfessor)
	if err != nil {
		t.Fatalf("expected find by role, got %v", err)
	}
	if len(usersByRole) != 1 {
		t.Fatalf("expected 1 professor, got %d", len(usersByRole))
	}
}

func TestUserRepositoryReturnsNotFound(t *testing.T) {
	setupTestDB(t)

	repo := usersInfrastructure.NewUserRepository()
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("expected automigrate, got %v", err)
	}

	_, err := repo.FindByID(999)
	if !errors.Is(err, usersDomain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	_, err = repo.FindByEmail("missing@example.com")
	if !errors.Is(err, usersDomain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
