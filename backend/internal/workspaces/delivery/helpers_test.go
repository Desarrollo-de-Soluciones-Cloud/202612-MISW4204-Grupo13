package delivery

import (
	"backend/internal/workspaces/domain"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestWorkspaceBindingHelpers(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(CreateWorkspaceRequest{})

	mapped := mapBindingErrors(err)
	if len(mapped) == 0 {
		t.Fatalf("expected workspace validation errors")
	}
}

func TestMapWorkspaceBindingErrorUnknownField(t *testing.T) {
	validate := validator.New()
	err := validate.Var("", "required")
	if err == nil {
		t.Fatalf("expected validation error")
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		t.Fatalf("expected validation errors")
	}
	if mapWorkspaceBindingError(validationErrors[0]) != nil {
		t.Fatalf("expected nil for unknown field")
	}
}

func TestIsWorkspaceValidationError(t *testing.T) {
	if !isWorkspaceValidationError(domain.ErrWorkspaceStateInvalid) {
		t.Fatalf("expected workspace state invalid to be validation error")
	}
	if isWorkspaceValidationError(errors.New("boom")) {
		t.Fatalf("did not expect generic error to be workspace validation error")
	}
}
