package infrastructure

import "testing"

func TestNewWeekRepository(t *testing.T) {
	repo := NewWeekRepository()
	if repo == nil {
		t.Fatalf("expected week repository")
	}
}
