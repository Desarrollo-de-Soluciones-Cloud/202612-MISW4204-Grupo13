package infrastructure

import "testing"

func TestNewTaskReader(t *testing.T) {
	reader := NewTaskReader()
	if reader == nil {
		t.Fatalf("expected task reader")
	}
}
