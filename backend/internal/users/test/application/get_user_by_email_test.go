package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"
)

func TestGetUserByEmailSuccess(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	getUserByEmail := applicationpkg.NewGetUserByEmail(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if err != nil {
		t.Fatalf("expected create user to succeed, got %v", err)
	}

	output, err := getUserByEmail.Execute(applicationpkg.GetUserByEmailInput{Email: testUserJohnEmailCaps})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Email != testUserJohnEmail {
		t.Fatalf("expected normalized email %q, got %q", testUserJohnEmail, output.Email)
	}
	if output.Password == "" {
		t.Fatal("expected password hash to be available for internal use")
	}
}

func TestGetUserByEmailInvalidEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	getUserByEmail := applicationpkg.NewGetUserByEmail(mockRepo)

	_, err := getUserByEmail.Execute(applicationpkg.GetUserByEmailInput{Email: "invalid-email"})
	if !errors.Is(err, domain.ErrUserEmailInvalid) {
		t.Fatalf("expected ErrUserEmailInvalid, got %v", err)
	}
}

func TestGetUserByEmailNotFound(t *testing.T) {
	mockRepo := newMockUserRepository()
	getUserByEmail := applicationpkg.NewGetUserByEmail(mockRepo)

	_, err := getUserByEmail.Execute(applicationpkg.GetUserByEmailInput{Email: "missing@example.com"})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
