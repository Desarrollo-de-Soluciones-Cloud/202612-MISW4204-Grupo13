package infrastructure_test

import (
	sharedDB "backend/internal/shared/database/testsupport"
	backendusers "backend/internal/users/domain"
	usersinfra "backend/internal/users/infrastructure"
	"errors"
	"testing"
)

func TestUserRepositorySQLiteCRUD(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &backendusers.User{})
	repo := usersinfra.NewUserRepository()

	user, err := backendusers.NewUser("Ana Gomez", "ana@example.com", "Password123", backendusers.RoleProfessor)
	if err != nil {
		t.Fatalf("expected valid user, got %v", err)
	}
	if err := repo.Create(user); err != nil {
		t.Fatalf("expected create user, got %v", err)
	}

	byID, err := repo.FindByID(user.ID)
	if err != nil || byID.Email != "ana@example.com" {
		t.Fatalf("expected find by id, got %v %#v", err, byID)
	}

	byEmail, err := repo.FindByEmail("ana@example.com")
	if err != nil || byEmail.ID != user.ID {
		t.Fatalf("expected find by email, got %v %#v", err, byEmail)
	}

	users, err := repo.FindAll()
	if err != nil || len(users) != 1 {
		t.Fatalf("expected 1 user, got %v %d", err, len(users))
	}

	byRole, err := repo.FindAllByRole(backendusers.RoleProfessor)
	if err != nil || len(byRole) != 1 {
		t.Fatalf("expected 1 professor, got %v %d", err, len(byRole))
	}

	user.Name = "Ana Maria"
	if err := repo.Update(user); err != nil {
		t.Fatalf("expected update, got %v", err)
	}
}

func TestUserRepositorySQLiteNotFound(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &backendusers.User{})
	repo := usersinfra.NewUserRepository()

	if _, err := repo.FindByID(999); !errors.Is(err, backendusers.ErrUserNotFound) {
		t.Fatalf("expected user not found, got %v", err)
	}
}
