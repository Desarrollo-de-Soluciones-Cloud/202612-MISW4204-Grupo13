package infrastructure_test

import (
	weeksInfrastructure "backend/internal/weeks/infrastructure"
	"testing"
)

func TestNewWeekRepository(t *testing.T) {
	repo := weeksInfrastructure.NewWeekRepository()
	if repo == nil {
		t.Fatalf("expected week repository")
	}
}
