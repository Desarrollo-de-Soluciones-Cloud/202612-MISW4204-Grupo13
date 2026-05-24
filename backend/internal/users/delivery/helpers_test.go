package delivery

import (
	"backend/internal/users/domain"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestUserMapBindingErrors(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(CreateUserRequest{})

	mapped := mapBindingErrors(err)
	if len(mapped) == 0 {
		t.Fatalf("expected mapped validation errors")
	}
	if !errors.Is(mapped[0], domain.ErrUserNameRequired) {
		t.Fatalf("expected first error to be user name required, got %v", mapped[0])
	}
}

func TestMapUserBindingErrorUnknownField(t *testing.T) {
	validate := validator.New()
	err := validate.Var("", "required")
	if err == nil {
		t.Fatalf("expected validation error")
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		t.Fatalf("expected validation errors")
	}
	if mapUserBindingError(validationErrors[0]) != nil {
		t.Fatalf("expected nil for unknown struct field")
	}
}
