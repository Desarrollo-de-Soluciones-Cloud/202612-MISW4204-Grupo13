package infrastructure

import "testing"

func TestNewAssignmentReader(t *testing.T) {
	reader := NewAssignmentReader()
	if reader == nil {
		t.Fatalf("expected assignment reader")
	}
}
