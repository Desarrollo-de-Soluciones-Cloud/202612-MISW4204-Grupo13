package domain

import (
	authDomain "backend/internal/auth/domain"
	"errors"
	"strings"
	"testing"
)

func TestValidateAuthEmail(t *testing.T) {
	if err := authDomain.ValidateAuthEmail(testAuthDomainEmail); err != nil {
		t.Fatalf("expected valid email, got %v", err)
	}
}

func TestValidateAuthEmailRejectsInvalidValue(t *testing.T) {
	err := authDomain.ValidateAuthEmail("invalid-email")
	if !errors.Is(err, authDomain.ErrAuthEmailInvalid) {
		t.Fatalf("expected ErrAuthEmailInvalid, got %v", err)
	}
}

func TestValidateAuthEmailRequired(t *testing.T) {
	err := authDomain.ValidateAuthEmail(" ")
	if !errors.Is(err, authDomain.ErrAuthEmailRequired) {
		t.Fatalf("expected ErrAuthEmailRequired, got %v", err)
	}
}

func TestValidateAuthPasswordRejectsShortPassword(t *testing.T) {
	err := authDomain.ValidateAuthPassword("short")
	if !errors.Is(err, authDomain.ErrAuthPasswordTooShort) {
		t.Fatalf("expected ErrAuthPasswordTooShort, got %v", err)
	}
}

func TestValidateAuthPasswordRequired(t *testing.T) {
	err := authDomain.ValidateAuthPassword(" ")
	if !errors.Is(err, authDomain.ErrAuthPasswordRequired) {
		t.Fatalf("expected ErrAuthPasswordRequired, got %v", err)
	}
}

func TestValidateAuthPasswordTooLong(t *testing.T) {
	err := authDomain.ValidateAuthPassword(strings.Repeat("a", 73))
	if !errors.Is(err, authDomain.ErrAuthPasswordTooLong) {
		t.Fatalf("expected ErrAuthPasswordTooLong, got %v", err)
	}
}

func TestValidateAuthPasswordValid(t *testing.T) {
	err := authDomain.ValidateAuthPassword(testAuthDomainPassword)
	if err != nil {
		t.Fatalf("expected valid password, got %v", err)
	}
}

func TestValidateTokenStringRejectsEmptyValue(t *testing.T) {
	err := authDomain.ValidateTokenString("")
	if !errors.Is(err, authDomain.ErrAuthTokenRequired) {
		t.Fatalf("expected ErrAuthTokenRequired, got %v", err)
	}
}
