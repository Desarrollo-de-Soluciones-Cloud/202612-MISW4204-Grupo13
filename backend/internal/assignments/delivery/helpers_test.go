package delivery

import (
	"backend/internal/assignments/domain"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestAssignmentMapBindingErrors(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(struct {
		UserID      uint   `validate:"required"`
		WorkspaceID uint   `validate:"required"`
		Role        string `validate:"required"`
		WeeklyHours int    `validate:"required,min=1"`
	}{})

	mapped := mapBindingErrors(err)
	if len(mapped) == 0 {
		t.Fatalf("expected mapped validation errors")
	}
	if !containsAssignmentError(mapped, domain.ErrAssignmentUserIDRequired) {
		t.Fatalf("expected user id required error, got %#v", mapped)
	}
	if !containsAssignmentError(mapped, domain.ErrAssignmentWorkspaceIDRequired) {
		t.Fatalf("expected workspace id required error, got %#v", mapped)
	}
}

func TestMapAssignmentBindingErrorUnknownField(t *testing.T) {
	validate := validator.New()
	err := validate.Var("", "required")
	if err == nil {
		t.Fatalf("expected validation error")
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		t.Fatalf("expected validation errors")
	}

	if mapAssignmentBindingError(validationErrors[0]) != nil {
		t.Fatalf("expected nil for unknown field")
	}
}

func TestIsAssignmentValidationError(t *testing.T) {
	if !isAssignmentValidationError(domain.ErrAssignmentRoleRequired) {
		t.Fatalf("expected assignment role required to be validation error")
	}
	if isAssignmentValidationError(errors.New("boom")) {
		t.Fatalf("did not expect generic error to be validation error")
	}
}

func containsAssignmentError(errs []error, target error) bool {
	for _, err := range errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
