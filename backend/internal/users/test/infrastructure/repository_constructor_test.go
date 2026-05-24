package infrastructure_test

import (
	usersInfrastructure "backend/internal/users/infrastructure"
	"testing"
)

func TestNewUserRepository(t *testing.T) {
	repo := usersInfrastructure.NewUserRepository()
	if repo == nil {
		t.Fatalf("expected user repository")
	}
}
