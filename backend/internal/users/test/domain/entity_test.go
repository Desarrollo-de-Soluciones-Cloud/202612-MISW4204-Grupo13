package domain

import (
	domainpkg "backend/internal/users/domain"
	"errors"
	"strings"
	"testing"
)

func TestNewUserSuccess(t *testing.T) {
	user, err := domainpkg.NewUser(" John Doe ", " John@Example.com ", "password123", domainpkg.RoleProfessor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != "John Doe" {
		t.Fatalf("expected normalized name, got %q", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Fatalf("expected normalized email, got %q", user.Email)
	}
	if user.GlobalRole != domainpkg.RoleProfessor {
		t.Fatalf("expected role %q, got %q", domainpkg.RoleProfessor, user.GlobalRole)
	}
	if user.Password != "password123" {
		t.Fatalf("expected password to remain unchanged in domain constructor, got %q", user.Password)
	}
}

func TestNewUserRejectsInvalidRole(t *testing.T) {
	_, err := domainpkg.NewUser("John Doe", "john@example.com", "password123", domainpkg.UserRole("visitor"))
	if err != domainpkg.ErrUserRoleInvalid {
		t.Fatalf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestUpdateProfileNormalizesValues(t *testing.T) {
	user, err := domainpkg.NewUser("John Doe", "john@example.com", "password123", domainpkg.RoleMonitor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = user.UpdateProfile(" Jane Doe ", " Jane@Example.com ", domainpkg.RoleAdmin)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != "Jane Doe" {
		t.Fatalf("expected normalized name, got %q", user.Name)
	}
	if user.Email != "jane@example.com" {
		t.Fatalf("expected normalized email, got %q", user.Email)
	}
	if user.GlobalRole != domainpkg.RoleAdmin {
		t.Fatalf("expected role %q, got %q", domainpkg.RoleAdmin, user.GlobalRole)
	}
}

func TestNewUserRejectsInvalidEmail(t *testing.T) {
	_, err := domainpkg.NewUser("John Doe", "invalid-email", "password123", domainpkg.RoleProfessor)
	if !errors.Is(err, domainpkg.ErrUserEmailInvalid) {
		t.Fatalf("expected ErrUserEmailInvalid, got %v", err)
	}
}

func TestNewUserRejectsShortPassword(t *testing.T) {
	_, err := domainpkg.NewUser("John Doe", "john@example.com", "short", domainpkg.RoleProfessor)
	if !errors.Is(err, domainpkg.ErrUserPasswordTooShort) {
		t.Fatalf("expected ErrUserPasswordTooShort, got %v", err)
	}
}

func TestUpdateProfileRejectsInvalidRole(t *testing.T) {
	user, err := domainpkg.NewUser("John Doe", "john@example.com", "password123", domainpkg.RoleProfessor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = user.UpdateProfile("John Doe", "john@example.com", domainpkg.UserRole("guest"))
	if !errors.Is(err, domainpkg.ErrUserRoleInvalid) {
		t.Fatalf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestNormalizeEmail(t *testing.T) {
	normalized := domainpkg.NormalizeEmail(" John.Doe@Example.com ")
	if normalized != "john.doe@example.com" {
		t.Fatalf("expected normalized email, got %q", normalized)
	}
}

func TestValidateUserNameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 101)
	err := domainpkg.ValidateUserName(longName)
	if !errors.Is(err, domainpkg.ErrUserNameTooLong) {
		t.Fatalf("expected ErrUserNameTooLong, got %v", err)
	}
}
