package infrastructure_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	reportsInfrastructure "backend/internal/reports/infrastructure"
	sharedDB "backend/internal/shared/database/testsupport"
	"testing"
)

func TestAssignmentReaderFindAllByWorkspaceIDSQLite(t *testing.T) {
	db := sharedDB.SetupSQLiteDB(t, &assignmentsDomain.Assignment{})

	dbAssignments := []assignmentsDomain.Assignment{
		{UserID: 10, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 6},
		{UserID: 11, WorkspaceID: 1, Role: assignmentsDomain.RoleMonitor, WeeklyHours: 4},
		{UserID: 12, WorkspaceID: 2, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 8},
	}
	for _, assignment := range dbAssignments {
		item := assignment
		if err := db.Create(&item).Error; err != nil {
			t.Fatalf("expected seed insert, got %v", err)
		}
	}

	reader := reportsInfrastructure.NewAssignmentReader()
	assignments, err := reader.FindAllByWorkspaceID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(assignments) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(assignments))
	}
	if assignments[0].WorkspaceID != 1 || assignments[1].WorkspaceID != 1 {
		t.Fatalf("expected only workspace 1 assignments, got %#v", assignments)
	}
}
