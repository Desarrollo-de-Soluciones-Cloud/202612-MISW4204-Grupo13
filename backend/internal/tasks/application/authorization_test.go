package application

import (
	usersDomain "backend/internal/users/domain"
	"testing"
)

func TestCanReadAllTasks(t *testing.T) {
	if !canReadAllTasks(usersDomain.RoleProfessor) {
		t.Fatalf("expected professor to read all tasks")
	}
	if !canReadAllTasks(usersDomain.RoleAdmin) {
		t.Fatalf("expected admin to read all tasks")
	}
	if canReadAllTasks(usersDomain.RoleMonitor) {
		t.Fatalf("did not expect monitor to read all tasks")
	}
}

func TestIsOperationalRole(t *testing.T) {
	if !isOperationalRole(usersDomain.RoleMonitor) {
		t.Fatalf("expected monitor to be operational role")
	}
	if !isOperationalRole(usersDomain.RoleAssistant) {
		t.Fatalf("expected assistant to be operational role")
	}
	if isOperationalRole(usersDomain.RoleProfessor) {
		t.Fatalf("did not expect professor to be operational role")
	}
}
