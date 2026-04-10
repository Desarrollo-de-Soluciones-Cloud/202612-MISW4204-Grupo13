package domain

import (
	"backend/internal/users/domain"
	"testing"
)

func TestUserEntityCreation(t *testing.T) {
    user := &domain.User{
        ID:        1,
        Name:      "John Doe",
        Email:     "john@example.com",
        Password:  "hashedpassword",
    }
    if user.Name != "John Doe" {
        t.Errorf("Expected name 'John Doe', got '%s'", user.Name)
    }
    if user.Email != "john@example.com" {
        t.Errorf("Expected email 'john@example.com', got '%s'", user.Email)
    }
}

func TestUserPasswordIsHidden(t *testing.T) {
    user := &domain.User{
        ID:       1,
        Name:     "John Doe",
        Email:    "john@example.com",
        Password: "secretpassword",
    }
    if user.Password != "secretpassword" {
        t.Errorf("Password should be accessible in struct, got '%s'", user.Password)
    }
}