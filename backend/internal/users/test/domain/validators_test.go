package domain

import (
	domainpkg "backend/internal/users/domain"
	"errors"
	"strings"
	"testing"
)

func TestValidateUserNameRequired(t *testing.T) {
	err := domainpkg.ValidateUserName("   ")
	if !errors.Is(err, domainpkg.ErrUserNameRequired) {
		t.Fatalf("expected ErrUserNameRequired, got %v", err)
	}
}

func TestValidateUserNameValid(t *testing.T) {
	err := domainpkg.ValidateUserName(testDomainUserJohnName)
	if err != nil {
		t.Fatalf("expected valid name, got %v", err)
	}
}

func TestValidateUserNameTooShort(t *testing.T) {
	err := domainpkg.ValidateUserName("Jo")
	if !errors.Is(err, domainpkg.ErrUserNameTooShort) {
		t.Fatalf("expected ErrUserNameTooShort, got %v", err)
	}
}

func TestValidateUserEmailRequired(t *testing.T) {
	err := domainpkg.ValidateUserEmail(" ")
	if !errors.Is(err, domainpkg.ErrUserEmailRequired) {
		t.Fatalf("expected ErrUserEmailRequired, got %v", err)
	}
}

func TestValidateUserPasswordRequired(t *testing.T) {
	err := domainpkg.ValidateUserPassword(" ")
	if !errors.Is(err, domainpkg.ErrUserPasswordRequired) {
		t.Fatalf("expected ErrUserPasswordRequired, got %v", err)
	}
}

func TestValidateUserPasswordTooLong(t *testing.T) {
	err := domainpkg.ValidateUserPassword(strings.Repeat("a", 73))
	if !errors.Is(err, domainpkg.ErrUserPasswordTooLong) {
		t.Fatalf("expected ErrUserPasswordTooLong, got %v", err)
	}
}

func TestValidateUserRoleRequired(t *testing.T) {
	err := domainpkg.ValidateUserRole("")
	if !errors.Is(err, domainpkg.ErrUserRoleRequired) {
		t.Fatalf("expected ErrUserRoleRequired, got %v", err)
	}
}

func TestValidateUserRoleValid(t *testing.T) {
	err := domainpkg.ValidateUserRole(domainpkg.RoleAssistant)
	if err != nil {
		t.Fatalf("expected valid role, got %v", err)
	}
}
