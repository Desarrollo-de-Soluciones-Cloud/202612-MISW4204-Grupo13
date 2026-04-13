package application

import (
	applicationpkg "backend/internal/assignments/application"
	"backend/internal/assignments/domain"
	"errors"
	"testing"
)

const (
	errCreateAssignmentMsg          = "expected no error creating assignment: %v"
	errCreateAssistantAssignmentMsg = "expected no error creating assistant assignment, got %v"
	errCreateMonitorToUpdateMsg     = "expected no error creating monitor assignment to update, got %v"
)

// TODO RF05: Reforzar RF05 con pruebas de delivery e integracion/E2E para respuestas HTTP bloqueantes en update admin.

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
		t.Fatalf(errCreateAssignmentMsg, err)
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
		t.Fatalf(errCreateAssignmentMsg, err)
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
		t.Fatalf(errCreateAssignmentMsg, err)
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

func TestUpdateAssignmentBlocksWhenAssistantHoursExceed22(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	first, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      100,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 12,
	})
	if err != nil {
		t.Fatalf("expected no error creating first assistant assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      100,
		WorkspaceID: 2,
		Role:        domain.RoleAssistant,
		WeeklyHours: 8,
	})
	if err != nil {
		t.Fatalf("expected no error creating second assistant assignment, got %v", err)
	}

	_, err = updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          first.ID,
		Role:        domain.RoleAssistant,
		WeeklyHours: 15,
	})
	if !errors.Is(err, domain.ErrAssignmentAssistantHoursLimitExceeded) {
		t.Fatalf("expected ErrAssignmentAssistantHoursLimitExceeded, got %v", err)
	}
}

func TestUpdateAssignmentBlocksWhenMonitorAssignmentsExceed3(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	for i := 1; i <= 3; i++ {
		_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
			UserID:      110,
			WorkspaceID: uint(i),
			Role:        domain.RoleMonitor,
			WeeklyHours: 2,
		})
		if err != nil {
			t.Fatalf("expected no error creating monitor assignment %d, got %v", i, err)
		}
	}

	toChange, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      110,
		WorkspaceID: 4,
		Role:        domain.RoleAssistant,
		WeeklyHours: 15,
	})
	if err != nil {
		t.Fatalf("expected no error creating assistant assignment to update, got %v", err)
	}

	_, err = updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          toChange.ID,
		Role:        domain.RoleMonitor,
		WeeklyHours: 1,
	})
	if !errors.Is(err, domain.ErrAssignmentMonitorCountLimitExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorCountLimitExceeded, got %v", err)
	}
}

func TestUpdateAssignmentBlocksWhenMonitorHoursExceed12(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	first, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      120,
		WorkspaceID: 1,
		Role:        domain.RoleMonitor,
		WeeklyHours: 6,
	})
	if err != nil {
		t.Fatalf("expected no error creating first monitor assignment, got %v", err)
	}

	_, err = createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      120,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 6,
	})
	if err != nil {
		t.Fatalf("expected no error creating second monitor assignment, got %v", err)
	}

	_, err = updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          first.ID,
		Role:        domain.RoleMonitor,
		WeeklyHours: 7,
	})
	if !errors.Is(err, domain.ErrAssignmentMonitorHoursLimitExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorHoursLimitExceeded, got %v", err)
	}
}

func TestUpdateAssignmentBlocksWhenMonitorExceedsFortyPercentOfAssistant(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      130,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf(errCreateAssistantAssignmentMsg, err)
	}

	toChange, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      130,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 3,
	})
	if err != nil {
		t.Fatalf(errCreateMonitorToUpdateMsg, err)
	}

	_, err = updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          toChange.ID,
		Role:        domain.RoleMonitor,
		WeeklyHours: 5,
	})
	if !errors.Is(err, domain.ErrAssignmentMonitorFortyPercentExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorFortyPercentExceeded, got %v", err)
	}
}

func TestUpdateAssignmentAllowsMonitorHoursAtRoundedFortyPercentLimit(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      140,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 11,
	})
	if err != nil {
		t.Fatalf(errCreateAssistantAssignmentMsg, err)
	}

	toChange, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      140,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 4,
	})
	if err != nil {
		t.Fatalf(errCreateMonitorToUpdateMsg, err)
	}

	output, err := updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          toChange.ID,
		Role:        domain.RoleMonitor,
		WeeklyHours: 5,
	})
	if err != nil {
		t.Fatalf("expected no error at rounded forty percent limit, got %v", err)
	}
	if output.WeeklyHours != 5 {
		t.Fatalf("expected weekly hours 5, got %d", output.WeeklyHours)
	}
}

func TestUpdateAssignmentValidCaseWithMixedRoles(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	_, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      150,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 20,
	})
	if err != nil {
		t.Fatalf(errCreateAssistantAssignmentMsg, err)
	}

	toChange, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      150,
		WorkspaceID: 2,
		Role:        domain.RoleMonitor,
		WeeklyHours: 5,
	})
	if err != nil {
		t.Fatalf(errCreateMonitorToUpdateMsg, err)
	}

	output, err := updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          toChange.ID,
		Role:        domain.RoleMonitor,
		WeeklyHours: 8,
	})
	if err != nil {
		t.Fatalf("expected no error updating valid mixed-role assignment, got %v", err)
	}
	if output.WeeklyHours != 8 {
		t.Fatalf("expected weekly hours 8, got %d", output.WeeklyHours)
	}
}

func TestUpdateAssignmentBlocksExactDuplicate(t *testing.T) {
	mockRepo := NewMockAssignmentRepository()
	createAssignment := applicationpkg.NewCreateAssignment(mockRepo)
	updateAssignment := applicationpkg.NewUpdateAssignment(mockRepo)

	first, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      160,
		WorkspaceID: 1,
		Role:        domain.RoleAssistant,
		WeeklyHours: 10,
	})
	if err != nil {
		t.Fatalf("expected no error creating first assignment, got %v", err)
	}

	second, err := createAssignment.Execute(applicationpkg.CreateAssignmentInput{
		UserID:      160,
		WorkspaceID: 1,
		Role:        domain.RoleMonitor,
		WeeklyHours: 4,
	})
	if err != nil {
		t.Fatalf("expected no error creating second assignment, got %v", err)
	}

	_, err = updateAssignment.Execute(applicationpkg.UpdateAssignmentInput{
		ID:          second.ID,
		Role:        domain.RoleAssistant,
		WeeklyHours: 4,
	})
	if !errors.Is(err, domain.ErrAssignmentAlreadyExists) {
		t.Fatalf("expected ErrAssignmentAlreadyExists, got %v", err)
	}

	if first.ID == second.ID {
		t.Fatalf("expected different assignment ids, got %d and %d", first.ID, second.ID)
	}
}
