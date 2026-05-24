package infrastructure_test

import (
	sharedDB "backend/internal/shared/database/testsupport"
	workspacesDomain "backend/internal/workspaces/domain"
	workspacesinfra "backend/internal/workspaces/infrastructure"
	"errors"
	"testing"
)

func TestWorkspaceRepositorySQLiteCRUD(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &workspacesDomain.Workspace{})
	repo := workspacesinfra.NewWorkspaceRepository()

	workspace, err := workspacesDomain.NewWorkspace(workspacesDomain.WorkspaceInput{
		PeriodID: 1, UserID: 10, Name: "Algorithms", Type: workspacesDomain.CourseType,
		InitialDate: "2026-06-02", FinalDate: "2026-06-30", Observations: "obs", State: workspacesDomain.ActiveState,
	})
	if err != nil {
		t.Fatalf("expected workspace, got %v", err)
	}
	if err := repo.Create(workspace); err != nil {
		t.Fatalf("expected create, got %v", err)
	}

	found, err := repo.FindByID(workspace.ID)
	if err != nil || found.Name != "Algorithms" {
		t.Fatalf("expected find by id, got %v %#v", err, found)
	}

	all, err := repo.FindAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("expected 1 workspace, got %v %d", err, len(all))
	}

	byPeriod, err := repo.FindByPeriodID(1)
	if err != nil || len(byPeriod) != 1 {
		t.Fatalf("expected 1 by period, got %v %d", err, len(byPeriod))
	}

	byUser, err := repo.FindByUserID(10)
	if err != nil || len(byUser) != 1 {
		t.Fatalf("expected 1 by user, got %v %d", err, len(byUser))
	}
}

func TestWorkspaceRepositorySQLiteDeleteNotFound(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &workspacesDomain.Workspace{})
	repo := workspacesinfra.NewWorkspaceRepository()

	if err := repo.Delete(999); !errors.Is(err, workspacesDomain.ErrWorkspaceNotFound) {
		t.Fatalf("expected workspace not found, got %v", err)
	}
}
