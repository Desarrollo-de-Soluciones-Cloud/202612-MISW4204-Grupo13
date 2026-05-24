package domain

import (
	reportsDomain "backend/internal/reports/domain"
	"errors"
	"testing"
)

func TestNewWeeklyReportSuccess(t *testing.T) {
	report, err := reportsDomain.NewWeeklyReport(1, 2, 3, 4, "reports/file.pdf")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.WorkspaceID != 1 || report.WeekID != 2 || report.AssignmentID != 3 || report.UserID != 4 {
		t.Fatalf("unexpected report: %+v", report)
	}
	if report.FilePath != "reports/file.pdf" {
		t.Fatalf("expected file path reports/file.pdf, got %q", report.FilePath)
	}
}

func TestNewWeeklyReportRejectsMissingFilePath(t *testing.T) {
	_, err := reportsDomain.NewWeeklyReport(1, 2, 3, 4, "")
	if !errors.Is(err, reportsDomain.ErrReportFilePathRequired) {
		t.Fatalf("expected ErrReportFilePathRequired, got %v", err)
	}
}

func TestNewWeeklyReportRejectsMissingAssignmentID(t *testing.T) {
	_, err := reportsDomain.NewWeeklyReport(1, 2, 0, 4, "reports/file.pdf")
	if !errors.Is(err, reportsDomain.ErrReportAssignmentIDRequired) {
		t.Fatalf("expected ErrReportAssignmentIDRequired, got %v", err)
	}
}
