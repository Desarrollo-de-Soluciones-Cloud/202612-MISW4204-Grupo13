package application

import (
	applicationpkg "backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	"errors"
	"testing"
)

func TestUpdateAssignmentSuccess(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	created, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 6,
	})
	if err != nil {
		t.Fatalf("expected no error creating assignment: %v", err)
	}

	output, err := updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          created.ID,
		Role:        domain.RoleAssistant,
		WeeklyHours: 12,
	})
	if err != nil {
		t.Fatalf("expected no error updating, got %v", err)
	}

	if output.Role != domain.RoleAssistant {
		t.Fatalf("expected role %q, got %q", domain.RoleAssistant, output.Role)
	}
	if output.WeeklyHours != 12 {
		t.Fatalf("expected weekly hours 12, got %d", output.WeeklyHours)
	}
	if output.UserID != created.UserID {
		t.Fatalf("expected user id %d to remain unchanged, got %d", created.UserID, output.UserID)
	}
	if output.WorkspaceID != created.WorkspaceID {
		t.Fatalf("expected workspace id %d to remain unchanged, got %d", created.WorkspaceID, output.WorkspaceID)
	}
}

func TestUpdateAssignmentNotFound(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	_, err := updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          999,
		Role:        domain.RoleMonitor,
		WeeklyHours: 8,
	})
	if !errors.Is(err, domain.ErrAssignmentNotFound) {
		t.Fatalf("expected ErrAssignmentNotFound, got %v", err)
	}
}

func TestUpdateAssignmentInvalidRole(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	created, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 6,
	})
	if err != nil {
		t.Fatalf("expected no error creating assignment: %v", err)
	}

	_, err = updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          created.ID,
		Role:        domain.AssignmentRole("invalid"),
		WeeklyHours: 8,
	})
	if !errors.Is(err, domain.ErrAssignmentRoleInvalid) {
		t.Fatalf("expected ErrAssignmentRoleInvalid, got %v", err)
	}
}

func TestUpdateAssignmentInvalidWeeklyHours(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	created, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 6,
	})
	if err != nil {
		t.Fatalf("expected no error creating assignment: %v", err)
	}

	_, err = updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          created.ID,
		Role:        domain.RoleAssistant,
		WeeklyHours: 0,
	})
	if !errors.Is(err, domain.ErrAssignmentWeeklyHoursInvalid) {
		t.Fatalf("expected ErrAssignmentWeeklyHoursInvalid, got %v", err)
	}
}
