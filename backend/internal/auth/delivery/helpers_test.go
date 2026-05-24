package delivery

import (
	"backend/internal/auth/domain"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestAuthMapBindingErrors(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(SignInRequest{})

	mapped := mapBindingErrors(err)
	if len(mapped) == 0 {
		t.Fatalf("expected mapped auth errors")
	}
	hasEmailRequired := false
	for _, mappedErr := range mapped {
		if errors.Is(mappedErr, domain.ErrAuthEmailRequired) {
			hasEmailRequired = true
			break
		}
	}
	if !hasEmailRequired {
		t.Fatalf("expected ErrAuthEmailRequired in mapped result, got %#v", mapped)
	}
}

func TestExtractBearerToken(t *testing.T) {
	token, err := extractBearerToken("Bearer abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "abc123" {
		t.Fatalf("expected token abc123, got %q", token)
	}

	if _, err := extractBearerToken(""); !errors.Is(err, domain.ErrAuthTokenRequired) {
		t.Fatalf("expected token required error, got %v", err)
	}
	if _, err := extractBearerToken("Basic token"); !errors.Is(err, domain.ErrAuthTokenInvalid) {
		t.Fatalf("expected token invalid error, got %v", err)
	}
}
