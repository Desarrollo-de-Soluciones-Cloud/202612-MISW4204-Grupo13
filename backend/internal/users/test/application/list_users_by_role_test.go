package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"
)

func TestListUsersByRoleSuccess(t *testing.T) {
	mockRepo := &MockListUsersRepository{
		users: []domain.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "hash", GlobalRole: domain.RoleProfessor},
			{ID: 2, Name: "Jane Doe", Email: "jane@example.com", Password: "hash", GlobalRole: domain.RoleAdmin},
			{ID: 3, Name: "Mike Doe", Email: "mike@example.com", Password: "hash", GlobalRole: domain.RoleProfessor},
		},
	}

	useCase := applicationpkg.NewListUsersByRole(mockRepo)
	output, err := useCase.Execute(applicationpkg.ListUsersByRoleInput{GlobalRole: domain.RoleProfessor})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(output.Users))
	}
	for _, user := range output.Users {
		if user.GlobalRole != domain.RoleProfessor {
			t.Fatalf("expected all users to have role %q, got %q", domain.RoleProfessor, user.GlobalRole)
		}
	}
}

func TestListUsersByRoleInvalidRole(t *testing.T) {
	mockRepo := &MockListUsersRepository{users: []domain.User{}}
	useCase := applicationpkg.NewListUsersByRole(mockRepo)

	_, err := useCase.Execute(applicationpkg.ListUsersByRoleInput{GlobalRole: domain.UserRole("guest")})
	if !errors.Is(err, domain.ErrUserRoleInvalid) {
		t.Fatalf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestListUsersByRoleEmptyResult(t *testing.T) {
	mockRepo := &MockListUsersRepository{
		users: []domain.User{
			{ID: 1, Name: "Jane Doe", Email: "jane@example.com", Password: "hash", GlobalRole: domain.RoleAdmin},
		},
	}

	useCase := applicationpkg.NewListUsersByRole(mockRepo)
	output, err := useCase.Execute(applicationpkg.ListUsersByRoleInput{GlobalRole: domain.RoleProfessor})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(output.Users))
	}
}
