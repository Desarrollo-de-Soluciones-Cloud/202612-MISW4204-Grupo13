package application

import (
	applicationpkg "backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	"errors"
	"testing"
)

func TestGetAssignmentByIDSuccess(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	getAssignment := applicationpkg.NewGetAssignmentByID(mockRepo)

	created, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      1,
		WorkspaceID: 3,
		Role:        domain.RoleMonitor,
		WeeklyHours: 8,
	})
	if err != nil {
		t.Fatalf("expected no error creating, got %v", err)
	}

	output, err := getAssignment.Execute(applicationpkg.GetAssignmentByIDInput{ID: created.ID})
	if err != nil {
		t.Fatalf("expected no error getting, got %v", err)
	}

	if output.ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, output.ID)
	}
	if output.UserID != 1 {
		t.Fatalf("expected user id 1, got %d", output.UserID)
	}
	if output.Role != domain.RoleMonitor {
		t.Fatalf("expected role %q, got %q", domain.RoleMonitor, output.Role)
	}
}

func TestGetAssignmentByIDNotFound(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	getAssignment := applicationpkg.NewGetAssignmentByID(mockRepo)

	_, err := getAssignment.Execute(applicationpkg.GetAssignmentByIDInput{ID: 999})
	if !errors.Is(err, domain.ErrAssignmentNotFound) {
		t.Fatalf("expected ErrAssignmentNotFound, got %v", err)
	}
}

func TestListAssignmentsByUserIDSuccess(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	listAssignments := applicationpkg.NewListAssignmentsByUserID(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID: 5, WorkspaceID: 1, Role: domain.RoleMonitor, WeeklyHours: 6,
	})
	if err != nil {
		t.Fatalf("expected no error creating first assignment: %v", err)
	}




	
	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID: 5, WorkspaceID: 2, Role: domain.RoleAssistant, WeeklyHours: 20,
	})
	if err != nil {
		t.Fatalf("expected no error creating second assignment: %v", err)
	}
	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID: 9, WorkspaceID: 2, Role: domain.RoleMonitor, WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf("expected no error creating third assignment: %v", err)
	}

	output, err := listAssignments.Execute(applicationpkg.ListAssignmentsByUserIDInput{UserID: 5})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Assignments) != 2 {
		t.Fatalf("expected 2 assignments for user 5, got %d", len(output.Assignments))
	}
}

func TestListAssignmentsByUserIDInvalidUserID(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	listAssignments := applicationpkg.NewListAssignmentsByUserID(mockRepo)

	_, err := listAssignments.Execute(applicationpkg.ListAssignmentsByUserIDInput{UserID: 0})
	if !errors.Is(err, domain.ErrAssignmentUserIDRequired) {
		t.Fatalf("expected ErrAssignmentUserIDRequired, got %v", err)
	}
}
