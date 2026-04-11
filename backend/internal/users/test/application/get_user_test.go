package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"
)

func TestGetUserByIDSuccess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	getUserByID := applicationpkg.NewGetUserByID(mockRepo)

	created, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Password:   "password123",
		GlobalRole: domain.RoleProfessor,
	})
	if err != nil {
		t.Fatalf("expected create user to succeed, got %v", err)
	}

	output, err := getUserByID.Execute(applicationpkg.GetUserByIDInput{ID: created.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Email != "john@example.com" {
		t.Fatalf("expected email 'john@example.com', got %q", output.Email)
	}
	if output.GlobalRole != domain.RoleProfessor {
		t.Fatalf("expected role %q, got %q", domain.RoleProfessor, output.GlobalRole)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	getUserByID := applicationpkg.NewGetUserByID(mockRepo)

	_, err := getUserByID.Execute(applicationpkg.GetUserByIDInput{ID: 999})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
