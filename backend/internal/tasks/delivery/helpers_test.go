package delivery

import (
	"backend/internal/tasks/domain"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestTaskMapBindingErrorsValidation(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(struct {
		AssignmentID uint   `validate:"required"`
		Title        string `validate:"required"`
		Description  string `validate:"required"`
		Status       string `validate:"required"`
		SpentHours   int    `validate:"required"`
		WeekStartDate string `validate:"required"`
	}{})

	mapped := mapBindingErrors(err)
	if len(mapped) == 0 {
		t.Fatalf("expected mapped errors")
	}
	if !containsTaskError(mapped, domain.ErrTaskAssignmentIDRequired) {
		t.Fatalf("expected assignment id required, got %#v", mapped)
	}
	if !containsTaskError(mapped, domain.ErrTaskWeekStartDateRequired) {
		t.Fatalf("expected week start date required, got %#v", mapped)
	}
}

func TestTaskMapBindingErrorsUnmarshalType(t *testing.T) {
	mapped := mapBindingErrors(&json.UnmarshalTypeError{Field: "spent_hours"})
	if len(mapped) != 1 || !errors.Is(mapped[0], domain.ErrTaskSpentHoursInvalid) {
		t.Fatalf("expected spent hours invalid, got %#v", mapped)
	}
}

func TestTaskMapBindingErrorsUnknownError(t *testing.T) {
	mapped := mapBindingErrors(errors.New("boom"))
	if len(mapped) != 1 || !errors.Is(mapped[0], domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %#v", mapped)
	}
}

func TestParseWeekStartDate(t *testing.T) {
	if _, err := parseWeekStartDate(""); !errors.Is(err, domain.ErrTaskWeekStartDateRequired) {
		t.Fatalf("expected required error, got %v", err)
	}
	if _, err := parseWeekStartDate("bad"); !errors.Is(err, domain.ErrTaskWeekStartDateInvalid) {
		t.Fatalf("expected invalid date error, got %v", err)
	}
	parsed, err := parseWeekStartDate("2026-04-06")
	if err != nil {
		t.Fatalf("expected valid date, got %v", err)
	}
	if parsed.Format(dateLayout) != "2026-04-06" {
		t.Fatalf("unexpected parsed date %v", parsed)
	}
}

func containsTaskError(errs []error, target error) bool {
	for _, err := range errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
