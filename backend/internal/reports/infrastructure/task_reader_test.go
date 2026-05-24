package infrastructure

import "testing"

func TestNewTaskReader(t *testing.T) {
	reader := NewTaskReader()
	if reader == nil {
		t.Fatalf("expected task reader")
	}
}

func TestTaskReaderFindAllByWorkspaceAndWeek(t *testing.T) {
	setupReportsDryRunDB(t)

	reader := NewTaskReader()
	tasks, err := reader.FindAllByWorkspaceAndWeek(3, 5, "2026-04-07")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tasks == nil {
		t.Fatalf("expected tasks slice")
	}
}
