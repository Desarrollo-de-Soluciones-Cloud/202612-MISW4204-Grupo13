package domain

import (
	reportsDomain "backend/internal/reports/domain"
	"errors"
	"testing"
)

const reportFilePath = "reports/file.pdf"

func TestNewWeeklyReportSuccess(t *testing.T) {
	report, err := reportsDomain.NewWeeklyReport(1, 2, 3, 4, reportFilePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.WorkspaceID != 1 || report.WeekID != 2 || report.AssignmentID != 3 || report.UserID != 4 {
		t.Fatalf("unexpected report: %+v", report)
	}
	if report.FilePath != reportFilePath {
		t.Fatalf("expected file path %s, got %q", reportFilePath, report.FilePath)
	}
}

func TestNewWeeklyReportRejectsMissingFilePath(t *testing.T) {
	_, err := reportsDomain.NewWeeklyReport(1, 2, 3, 4, "")
	if !errors.Is(err, reportsDomain.ErrReportFilePathRequired) {
		t.Fatalf("expected ErrReportFilePathRequired, got %v", err)
	}
}

func TestNewWeeklyReportRejectsMissingAssignmentID(t *testing.T) {
	_, err := reportsDomain.NewWeeklyReport(1, 2, 0, 4, reportFilePath)
	if !errors.Is(err, reportsDomain.ErrReportAssignmentIDRequired) {
		t.Fatalf("expected ErrReportAssignmentIDRequired, got %v", err)
	}
}
