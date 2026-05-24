package domain_test

import (
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"testing"
)

func TestValidateWorkspaceName(t *testing.T) {
	if !errors.Is(workspacesDomain.ValidateWorkspaceName(""), workspacesDomain.ErrWorkspaceNameRequired) {
		t.Fatalf("expected ErrWorkspaceNameRequired")
	}
	if err := workspacesDomain.ValidateWorkspaceName("Algorithms"); err != nil {
		t.Fatalf("expected valid workspace name, got %v", err)
	}
}

func TestValidateWorkspaceTypeAndState(t *testing.T) {
	if !errors.Is(workspacesDomain.ValidateWorkspaceType(""), workspacesDomain.ErrWorkspaceTypeRequired) {
		t.Fatalf("expected ErrWorkspaceTypeRequired")
	}
	if !errors.Is(workspacesDomain.ValidateWorkspaceType("invalid"), workspacesDomain.ErrWorkspaceTypeInvalid) {
		t.Fatalf("expected ErrWorkspaceTypeInvalid")
	}
	if !errors.Is(workspacesDomain.ValidateWorkspaceState(""), workspacesDomain.ErrWorkspaceStateRequired) {
		t.Fatalf("expected ErrWorkspaceStateRequired")
	}
	if err := workspacesDomain.ValidateWorkspaceState(workspacesDomain.ActiveState); err != nil {
		t.Fatalf("expected valid workspace state, got %v", err)
	}
}

func TestValidateWorkspaceDates(t *testing.T) {
	if !errors.Is(workspacesDomain.ValidateWorkspaceInitialDate("bad-date"), workspacesDomain.ErrWorkspaceInitialDateWrongFormat) {
		t.Fatalf("expected initial date format error")
	}
	if !errors.Is(workspacesDomain.ValidateWorkspaceFinalDate("bad-date"), workspacesDomain.ErrWorkspaceFinalDateWrongFormat) {
		t.Fatalf("expected final date format error")
	}
	if !errors.Is(workspacesDomain.ValidateWorkspaceDateSequence("2026-06-20", "2026-06-10"), workspacesDomain.ErrWorkspaceDateSequenceInvalid) {
		t.Fatalf("expected date sequence error")
	}
	if err := workspacesDomain.ValidateWorkspaceDateSequence("2026-06-10", "2026-06-20"); err != nil {
		t.Fatalf("expected valid date sequence, got %v", err)
	}
}
