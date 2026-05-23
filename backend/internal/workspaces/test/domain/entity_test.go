package domain_test

import (
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
)

const (
	testDomainWorkspaceStartDate = "2026-06-01"
	testDomainWorkspaceEndDate   = "2026-06-30"
)

func TestNewWorkspaceSuccess(t *testing.T) {
	workspace, err := workspacesDomain.NewWorkspace(workspacesDomain.WorkspaceInput{
		PeriodID:     1,
		UserID:       2,
		Name:         " Algorithms ",
		Type:         workspacesDomain.WorkspaceType("course"),
		InitialDate:  testDomainWorkspaceStartDate,
		FinalDate:    testDomainWorkspaceEndDate,
		Observations: "notes",
		State:        workspacesDomain.ActiveState,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if workspace.Name != "Algorithms" {
		t.Fatalf("expected normalized name, got %q", workspace.Name)
	}
}

func TestNewWorkspaceRejectsMissingProfessorID(t *testing.T) {
	_, err := workspacesDomain.NewWorkspace(workspacesDomain.WorkspaceInput{
		PeriodID:     1,
		UserID:       0,
		Name:         "Algorithms",
		Type:         workspacesDomain.WorkspaceType("course"),
		InitialDate:  testDomainWorkspaceStartDate,
		FinalDate:    testDomainWorkspaceEndDate,
		Observations: "notes",
		State:        workspacesDomain.ActiveState,
	})
	if !errors.Is(err, workspacesDomain.ErrWorkspaceUserIDRequired) {
		t.Fatalf("expected ErrWorkspaceUserIDRequired, got %v", err)
	}
}

func TestUpdateWorkspaceRejectsInvalidDateSequence(t *testing.T) {
	workspace, err := workspacesDomain.NewWorkspace(workspacesDomain.WorkspaceInput{
		PeriodID:     1,
		UserID:       2,
		Name:         "Algorithms",
		Type:         workspacesDomain.WorkspaceType("course"),
		InitialDate:  testDomainWorkspaceStartDate,
		FinalDate:    testDomainWorkspaceEndDate,
		Observations: "notes",
		State:        workspacesDomain.ActiveState,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = workspace.UpdateWorkspace(
		"Algorithms II",
		workspacesDomain.WorkspaceType("course"),
		"2026-06-20",
		"2026-06-10",
		"notes",
		workspacesDomain.ActiveState,
	)
	if !errors.Is(err, workspacesDomain.ErrWorkspaceDateSequenceInvalid) {
		t.Fatalf("expected ErrWorkspaceDateSequenceInvalid, got %v", err)
	}
}
