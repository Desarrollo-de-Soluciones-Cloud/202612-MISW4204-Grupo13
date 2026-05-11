package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCreateUserSuccess(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	input := applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmailCaps,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	}

	output, err := createUser.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != testUserJohnName {
		t.Errorf("expected name %q, got %q", testUserJohnName, output.Name)
	}
	if output.Email != testUserJohnEmail {
		t.Errorf("expected normalized email %q, got %q", testUserJohnEmail, output.Email)
	}
	if output.GlobalRole != domain.RoleProfessor {
		t.Errorf("expected role %q, got %q", domain.RoleProfessor, output.GlobalRole)
	}

	storedUser, err := mockRepo.FindByID(output.ID)
	if err != nil {
		t.Fatalf("expected stored user, got %v", err)
	}
	if storedUser.Password == testUserPassword123 {
		t.Fatal("expected stored password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(testUserPassword123)); err != nil {
		t.Fatalf("expected password hash to match original password: %v", err)
	}
}

func TestCreateUserInvalidName(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	input := applicationpkg.CreateUserInput{
		Name:       "",
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	}

	_, err := createUser.Execute(input)
	if !errors.Is(err, domain.ErrUserNameRequired) {
		t.Errorf("expected ErrUserNameRequired, got %v", err)
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	input := applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	}

	_, err := createUser.Execute(input)
	if err != nil {
		t.Fatalf("first creation should succeed, got %v", err)
	}

	input2 := applicationpkg.CreateUserInput{
		Name:       testUserJaneName,
		Email:      testUserJohnEmailCaps,
		Password:   "password456",
		GlobalRole: domain.RoleAdmin,
	}

	_, err = createUser.Execute(input2)
	if !errors.Is(err, domain.ErrUserEmailAlreadyInUse) {
		t.Errorf("expected ErrUserEmailAlreadyInUse, got %v", err)
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      "invalid-email",
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserEmailInvalid) {
		t.Errorf("expected ErrUserEmailInvalid, got %v", err)
	}
}

func TestCreateUserInvalidRole(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.UserRole("guest"),
	})
	if !errors.Is(err, domain.ErrUserRoleInvalid) {
		t.Errorf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestCreateUserPasswordTooShort(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   "short",
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserPasswordTooShort) {
		t.Errorf("expected ErrUserPasswordTooShort, got %v", err)
	}
}

func TestCreateUserPropagatesUnexpectedFindByEmailError(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockRepo.findByEmailErr = errors.New("database failure")
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if err == nil || err.Error() != "database failure" {
		t.Fatalf("expected database failure, got %v", err)
	}
}

func TestCreateUserPropagatesCreateError(t *testing.T) {
	mockRepo := newMockUserRepository()
	mockRepo.createErr = errors.New("create failure")
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   testUserPassword123,
		GlobalRole: domain.RoleProfessor,
	})
	if err == nil || err.Error() != "create failure" {
		t.Fatalf("expected create failure, got %v", err)
	}
}

func TestCreateUserRejectsRequiredPassword(t *testing.T) {
	mockRepo := newMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       testUserJohnName,
		Email:      testUserJohnEmail,
		Password:   "   ",
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserPasswordRequired) {
		t.Fatalf("expected ErrUserPasswordRequired, got %v", err)
	}
}
