package domain

import (
	authDomain "backend/internal/auth/domain"
	"errors"
	"testing"
)

func TestValidateAuthEmail(t *testing.T) {
	if err := authDomain.ValidateAuthEmail("john@example.com"); err != nil {
		t.Fatalf("expected valid email, got %v", err)
	}
}

func TestValidateAuthEmailRejectsInvalidValue(t *testing.T) {
	err := authDomain.ValidateAuthEmail("invalid-email")
	if !errors.Is(err, authDomain.ErrAuthEmailInvalid) {
		t.Fatalf("expected ErrAuthEmailInvalid, got %v", err)
	}
}

func TestValidateAuthPasswordRejectsShortPassword(t *testing.T) {
	err := authDomain.ValidateAuthPassword("short")
	if !errors.Is(err, authDomain.ErrAuthPasswordTooShort) {
		t.Fatalf("expected ErrAuthPasswordTooShort, got %v", err)
	}
}

func TestValidateTokenStringRejectsEmptyValue(t *testing.T) {
	err := authDomain.ValidateTokenString("")
	if !errors.Is(err, authDomain.ErrAuthTokenRequired) {
		t.Fatalf("expected ErrAuthTokenRequired, got %v", err)
	}
}
