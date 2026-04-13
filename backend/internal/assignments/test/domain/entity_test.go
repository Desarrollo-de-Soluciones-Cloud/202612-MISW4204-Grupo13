package domain

import (
	domainpkg "backend/internal/assignments/domain"
	"errors"
	"testing"
)

func TestNewAssignmentSuccess(t *testing.T) {
	assignment, err := domainpkg.NewAssignment(1, 2, domainpkg.RoleMonitor, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if assignment.UserID != 1 {
		t.Fatalf("expected user id 1, got %d", assignment.UserID)
	}
	if assignment.WorkspaceID != 2 {
		t.Fatalf("expected workspace id 2, got %d", assignment.WorkspaceID)
	}
	if assignment.Role != domainpkg.RoleMonitor {
		t.Fatalf("expected role %q, got %q", domainpkg.RoleMonitor, assignment.Role)
	}
	if assignment.WeeklyHours != 10 {
		t.Fatalf("expected weekly hours 10, got %d", assignment.WeeklyHours)
	}
}

func TestNewAssignmentInvalidRole(t *testing.T) {
	_, err := domainpkg.NewAssignment(1, 2, domainpkg.AssignmentRole("invalid"), 10)
	if !errors.Is(err, domainpkg.ErrAssignmentRoleInvalid) {
		t.Fatalf("expected ErrAssignmentRoleInvalid, got %v", err)
	}
}

func TestNewAssignmentMissingUserID(t *testing.T) {
	_, err := domainpkg.NewAssignment(0, 2, domainpkg.RoleAssistant, 10)
	if !errors.Is(err, domainpkg.ErrAssignmentUserIDRequired) {
		t.Fatalf("expected ErrAssignmentUserIDRequired, got %v", err)
	}
}

func TestNewAssignmentMissingWorkspaceID(t *testing.T) {
	_, err := domainpkg.NewAssignment(1, 0, domainpkg.RoleAssistant, 10)
	if !errors.Is(err, domainpkg.ErrAssignmentWorkspaceIDRequired) {
		t.Fatalf("expected ErrAssignmentWorkspaceIDRequired, got %v", err)
	}
}

func TestNewAssignmentInvalidWeeklyHours(t *testing.T) {
	_, err := domainpkg.NewAssignment(1, 2, domainpkg.RoleAssistant, 0)
	if !errors.Is(err, domainpkg.ErrAssignmentWeeklyHoursInvalid) {
		t.Fatalf("expected ErrAssignmentWeeklyHoursInvalid, got %v", err)
	}
}
