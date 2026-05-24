package infrastructure

import "testing"

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository()
	if repo == nil {
		t.Fatalf("expected user repository")
	}
}
