package infrastructure_test

import (
	assignmentsInfrastructure "backend/internal/assignments/infrastructure"
	"testing"
)

func TestNewAssignmentRepository(t *testing.T) {
	repo := assignmentsInfrastructure.NewAssignmentRepository()
	if repo == nil {
		t.Fatalf("expected assignment repository")
	}
}
