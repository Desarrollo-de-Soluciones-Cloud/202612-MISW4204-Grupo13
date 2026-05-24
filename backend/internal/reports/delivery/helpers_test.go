package delivery

import (
	reportsDomain "backend/internal/reports/domain"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestMapBindingErrorsForGenerateWeeklyReportsRequest(t *testing.T) {
	validate := validator.New()
	err := validate.Struct(GenerateWeeklyReportsRequest{})

	mapped := mapBindingErrors(err)
	if len(mapped) != 2 {
		t.Fatalf("expected 2 mapped errors, got %d", len(mapped))
	}
	if !errors.Is(mapped[0], reportsDomain.ErrReportWorkspaceIDRequired) && !errors.Is(mapped[1], reportsDomain.ErrReportWorkspaceIDRequired) {
		t.Fatalf("expected workspace id required error in mapped result")
	}
}

func TestMapBindingErrorsHandlesUnmarshalTypeError(t *testing.T) {
	mapped := mapBindingErrors(&json.UnmarshalTypeError{})
	if len(mapped) != 1 || !errors.Is(mapped[0], reportsDomain.ErrReportInvalidInput) {
		t.Fatalf("expected invalid input error, got %#v", mapped)
	}
}

func TestReportHelperParsers(t *testing.T) {
	if _, err := parseRequiredWorkspaceID(""); !errors.Is(err, reportsDomain.ErrReportWorkspaceFilterRequired) {
		t.Fatalf("expected missing workspace filter error, got %v", err)
	}

	value, err := parseOptionalResourceID("12")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value == nil || *value != 12 {
		t.Fatalf("expected parsed id 12, got %#v", value)
	}
}

func TestIsReportValidationError(t *testing.T) {
	if !isReportValidationError(reportsDomain.ErrReportFilePathRequired) {
		t.Fatalf("expected report file path required to be validation error")
	}
	if isReportValidationError(errors.New("boom")) {
		t.Fatalf("did not expect generic error to be report validation error")
	}
}
