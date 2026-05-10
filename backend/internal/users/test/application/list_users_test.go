package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"testing"
)

func TestListUsersEmpty(t *testing.T) {
	mockRepo := newMockUserRepository()
	listUsers := applicationpkg.NewListUsers(mockRepo)
	output, err := listUsers.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Users) != 0 {
		t.Errorf("expected 0 users, got %d", len(output.Users))
	}
}

func TestListUsersWithData(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockRepo.byID[1] = &domain.User{ID: 1, Name: testUserJohnName, Email: testUserJohnEmail, Password: "password-hash", GlobalRole: domain.RoleProfessor}
	mockRepo.byID[2] = &domain.User{ID: 2, Name: testUserJaneName, Email: testUserJaneEmail, Password: "password-hash", GlobalRole: domain.RoleAdmin}
	mockRepo.users[testUserJohnEmail] = mockRepo.byID[1]
	mockRepo.users[testUserJaneEmail] = mockRepo.byID[2]

	listUsers := applicationpkg.NewListUsers(mockRepo)
	output, err := listUsers.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(output.Users))
	}
	if output.Users[0].GlobalRole != domain.RoleProfessor {
		t.Errorf("expected first user role %q, got %q", domain.RoleProfessor, output.Users[0].GlobalRole)
	}
	if output.Users[1].GlobalRole != domain.RoleAdmin {
		t.Errorf("expected second user role %q, got %q", domain.RoleAdmin, output.Users[1].GlobalRole)
	}
}
