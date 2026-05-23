package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"
)

const testCreateUserSuccessMsg = "expected create user to succeed, got %v"

func TestUpdateUserSuccess(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	updateUser := applicationpkg.NewUpdateUser(mockRepo)

	created, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if err != nil {
		t.Fatalf(testCreateUserSuccessMsg, err)
	}

	output, err := updateUser.Execute(applicationpkg.UpdateUserInput{
		ID:         created.ID,
		Name:       testUserJaneName,
		Email:      testUserJaneEmailCaps,
		GlobalRole: domain.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != testUserJaneName {
		t.Fatalf("expected updated name %q, got %q", testUserJaneName, output.Name)
	}
	if output.Email != testUserJaneEmail {
		t.Fatalf("expected normalized email %q, got %q", testUserJaneEmail, output.Email)
	}
	if output.GlobalRole != domain.RoleAdmin {
		t.Fatalf("expected role %q, got %q", domain.RoleAdmin, output.GlobalRole)
	}
}

func TestUpdateUserRejectsDuplicateEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	updateUser := applicationpkg.NewUpdateUser(mockRepo)

	first, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if err != nil {
		t.Fatalf("expected first create user to succeed, got %v", err)
	}

	_, err = createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJaneName,
		Email:      testUserJaneEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("expected second create user to succeed, got %v", err)
	}

	_, err = updateUser.Execute(applicationpkg.UpdateUserInput{
		ID:         first.ID,
		Name:       testUserJohnName,
		Email:      testUserJaneEmail,
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserEmailAlreadyInUse) {
		t.Fatalf("expected ErrUserEmailAlreadyInUse, got %v", err)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	mockRepo := newMockUserRepository()
	updateUser := applicationpkg.NewUpdateUser(mockRepo)

	_, err := updateUser.Execute(applicationpkg.UpdateUserInput{
		ID:         999,
		Name:       testUserJaneName,
		Email:      testUserJaneEmail,
		GlobalRole: domain.RoleAdmin,
	})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdateUserInvalidRole(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	updateUser := applicationpkg.NewUpdateUser(mockRepo)

	created, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if err != nil {
		t.Fatalf(testCreateUserSuccessMsg, err)
	}

	_, err = updateUser.Execute(applicationpkg.UpdateUserInput{
		ID:         created.ID,
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		GlobalRole: domain.UserRole("guest"),
	})
	if !errors.Is(err, domain.ErrUserRoleInvalid) {
		t.Fatalf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestUpdateUserPreservesPasswordAndReplacesEmailLookup(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	updateUser := applicationpkg.NewUpdateUser(mockRepo)

	created, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if err != nil {
		t.Fatalf(testCreateUserSuccessMsg, err)
	}

	originalUser, err := mockRepo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected stored user, got %v", err)
	}
	originalPasswordHash := originalUser.Password

	_, err = updateUser.Execute(applicationpkg.UpdateUserInput{
		ID:         created.ID,
		Name:       testUserJaneName,
		Email:      testUserJaneEmail,
		GlobalRole: domain.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	updatedUser, err := mockRepo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected updated user, got %v", err)
	}
	if updatedUser.Password != originalPasswordHash {
		t.Fatalf("expected password hash to remain unchanged")
	}

	_, err = mockRepo.FindByEmail(testUserJohnEmail)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected old email lookup to fail with ErrUserNotFound, got %v", err)
	}

	foundByNewEmail, err := mockRepo.FindByEmail(testUserJaneEmail)
	if err != nil {
		t.Fatalf("expected new email lookup to succeed, got %v", err)
	}
	if foundByNewEmail.ID != created.ID {
		t.Fatalf("expected new email to map to user ID %d, got %d", created.ID, foundByNewEmail.ID)
	}
}
