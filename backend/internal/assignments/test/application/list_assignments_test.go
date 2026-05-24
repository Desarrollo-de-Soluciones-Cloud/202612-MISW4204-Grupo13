package application

import (
	"backend/internal/assignments/domain"
	"testing"
)

func TestListAllAssignmentsReturnsMappedAssignments(t *testing.T) {
	repo := newMockAssignmentRepository()
	repo.byID[1] = &domain.Assignment{ID: 1, UserID: 2, WorkspaceID: 3, Role: domain.RoleMonitor, WeeklyHours: 8}

	output, err := NewListAllAssignments(repo).Execute(ListAllAssignmentsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(output.Assignments))
	}
	if output.Assignments[0].WorkspaceID != 3 {
		t.Fatalf("expected workspace id 3, got %d", output.Assignments[0].WorkspaceID)
	}
}

func TestListAssignmentsByWorkspaceRejectsInvalidProfessorID(t *testing.T) {
	repo := newMockAssignmentRepository()

	_, err := NewListAssignmentsByWorkspace(repo).Execute(ListAssignmentsByWorkspaceInput{})
	if err == nil {
		t.Fatalf("expected validation error for empty professor id")
	}
}

func TestListAssignmentsByWorkspaceReturnsRepositoryResults(t *testing.T) {
	repo := newMockAssignmentRepository()
	repo.byID[1] = &domain.Assignment{ID: 1, UserID: 99, WorkspaceID: 7, Role: domain.RoleAssistant, WeeklyHours: 16}

	output, err := NewListAssignmentsByWorkspace(repo).Execute(ListAssignmentsByWorkspaceInput{ProfessorID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatalf("expected output to be present")
	}
}
