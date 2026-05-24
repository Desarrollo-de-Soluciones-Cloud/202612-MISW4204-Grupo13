package infrastructure_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	sharedDB "backend/internal/shared/database/testsupport"
	assignmentsinfra "backend/internal/assignments/infrastructure"
	"errors"
	"testing"
)

func TestAssignmentRepositorySQLiteCRUD(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &assignmentsDomain.Assignment{})
	repo := assignmentsinfra.NewAssignmentRepository()

	assignment, err := assignmentsDomain.NewAssignment(2, 7, assignmentsDomain.RoleMonitor, 8)
	if err != nil {
		t.Fatalf("expected assignment, got %v", err)
	}
	if err := repo.Create(assignment); err != nil {
		t.Fatalf("expected create, got %v", err)
	}

	found, err := repo.FindByID(assignment.ID)
	if err != nil || found.UserID != 2 {
		t.Fatalf("expected find by id, got %v %#v", err, found)
	}

	all, err := repo.FindAll()
	if err != nil || len(all) != 1 {
		t.Fatalf("expected 1 assignment, got %v %d", err, len(all))
	}

	byUser, err := repo.FindAllByUserID(2)
	if err != nil || len(byUser) != 1 {
		t.Fatalf("expected 1 by user, got %v %d", err, len(byUser))
	}

	sum, err := repo.SumWeeklyHoursByUserAndRole(2, assignmentsDomain.RoleMonitor)
	if err != nil || sum != 8 {
		t.Fatalf("expected sum 8, got %v %d", err, sum)
	}

	count, err := repo.CountAssignmentsByUserAndRole(2, assignmentsDomain.RoleMonitor)
	if err != nil || count != 1 {
		t.Fatalf("expected count 1, got %v %d", err, count)
	}
}

func TestAssignmentRepositorySQLiteNotFound(t *testing.T) {
	sharedDB.SetupSQLiteDB(t, &assignmentsDomain.Assignment{})
	repo := assignmentsinfra.NewAssignmentRepository()

	if _, err := repo.FindByID(999); !errors.Is(err, assignmentsDomain.ErrAssignmentNotFound) {
		t.Fatalf("expected assignment not found, got %v", err)
	}
}
