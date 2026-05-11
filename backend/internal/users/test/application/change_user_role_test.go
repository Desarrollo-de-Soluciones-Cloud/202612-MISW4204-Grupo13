package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"
)

func TestChangeUserRoleSuccess(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	changeUserRole := applicationpkg.NewChangeUserRole(mockRepo)

	created, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleMonitor,
	})
	if err != nil {
		t.Fatalf("expected create user to succeed, got %v", err)
	}

	output, err := changeUserRole.Execute(applicationpkg.ChangeUserRoleInput{
		ID:         created.ID,
		GlobalRole: domain.RoleProfessor,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.GlobalRole != domain.RoleProfessor {
		t.Fatalf("expected role %q, got %q", domain.RoleProfessor, output.GlobalRole)
	}

	storedUser, err := mockRepo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected stored user, got %v", err)
	}
	if storedUser.GlobalRole != domain.RoleProfessor {
		t.Fatalf("expected persisted role %q, got %q", domain.RoleProfessor, storedUser.GlobalRole)
	}
	if storedUser.Name != testUserJohnName {
		t.Fatalf("expected name to remain unchanged, got %q", storedUser.Name)
	}
	if storedUser.Email != testUserJohnEmail {
		t.Fatalf("expected email to remain unchanged, got %q", storedUser.Email)
	}
}

func TestChangeUserRoleInvalidRole(t *testing.T) {
	mockRepo := newMockUserRepository()
	changeUserRole := applicationpkg.NewChangeUserRole(mockRepo)

	_, err := changeUserRole.Execute(applicationpkg.ChangeUserRoleInput{
		ID:         1,
		GlobalRole: domain.UserRole("guest"),
	})
	if !errors.Is(err, domain.ErrUserRoleInvalid) {
		t.Fatalf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestChangeUserRoleRequiresRole(t *testing.T) {
	mockRepo := newMockUserRepository()
	changeUserRole := applicationpkg.NewChangeUserRole(mockRepo)

	_, err := changeUserRole.Execute(applicationpkg.ChangeUserRoleInput{
		ID:         1,
		GlobalRole: "",
	})
	if !errors.Is(err, domain.ErrUserRoleRequired) {
		t.Fatalf("expected ErrUserRoleRequired, got %v", err)
	}
}

func TestChangeUserRoleUserNotFound(t *testing.T) {
	mockRepo := newMockUserRepository()
	changeUserRole := applicationpkg.NewChangeUserRole(mockRepo)

	_, err := changeUserRole.Execute(applicationpkg.ChangeUserRoleInput{
		ID:         999,
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestChangeUserRolePropagatesUpdateError(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	changeUserRole := applicationpkg.NewChangeUserRole(mockRepo)

	created, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleMonitor,
	})
	if err != nil {
		t.Fatalf("expected create user to succeed, got %v", err)
	}

	mockRepo.updateErr = errors.New("update failure")

	_, err = changeUserRole.Execute(applicationpkg.ChangeUserRoleInput{
		ID:         created.ID,
		GlobalRole: domain.RoleProfessor,
	})
	if err == nil || err.Error() != "update failure" {
		t.Fatalf("expected update failure, got %v", err)
	}
}
