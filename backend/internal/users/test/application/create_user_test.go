package application

import (
	applicationpkg "backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	users          map[string]*domain.User
	byID           map[uint]*domain.User
	nextID         uint
	createErr      error
	updateErr      error
	findByEmailErr error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*domain.User),
		byID:   make(map[uint]*domain.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) FindAll() ([]domain.User, error) {
	users := make([]domain.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, nil
}

func (m *MockUserRepository) FindAllByRole(role domain.UserRole) ([]domain.User, error) {
	users := make([]domain.User, 0)
	for _, u := range m.users {
		if u.GlobalRole == role {
			users = append(users, *u)
		}
	}
	return users, nil
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	if user, ok := m.byID[id]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) Update(user *domain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.byID[user.ID]; !ok {
		return domain.ErrUserNotFound
	}

	for email, existingUser := range m.users {
		if existingUser.ID == user.ID && email != user.Email {
			delete(m.users, email)
		}
	}

	m.byID[user.ID] = user
	m.users[user.Email] = user
	return nil
}

func TestCreateUserSuccess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	input := applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "JOHN@example.com",
		Password:   "password123",
		GlobalRole: domain.RoleProfessor,
	}

	output, err := createUser.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got %q", output.Name)
	}
	if output.Email != "john@example.com" {
		t.Errorf("expected normalized email 'john@example.com', got %q", output.Email)
	}
	if output.GlobalRole != domain.RoleProfessor {
		t.Errorf("expected role %q, got %q", domain.RoleProfessor, output.GlobalRole)
	}

	storedUser, err := mockRepo.FindByID(output.ID)
	if err != nil {
		t.Fatalf("expected stored user, got %v", err)
	}
	if storedUser.Password == "password123" {
		t.Fatal("expected stored password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte("password123")); err != nil {
		t.Fatalf("expected password hash to match original password: %v", err)
	}
}

func TestCreateUserInvalidName(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	input := applicationpkg.CreateUserInput{
		Name:       "",
		Email:      "john@example.com",
		Password:   "password123",
		GlobalRole: domain.RoleProfessor,
	}

	_, err := createUser.Execute(input)
	if !errors.Is(err, domain.ErrUserNameRequired) {
		t.Errorf("expected ErrUserNameRequired, got %v", err)
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)
	input := applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Password:   "password123",
		GlobalRole: domain.RoleProfessor,
	}

	_, err := createUser.Execute(input)
	if err != nil {
		t.Fatalf("first creation should succeed, got %v", err)
	}

	input2 := applicationpkg.CreateUserInput{
		Name:       "Jane Doe",
		Email:      "JOHN@example.com",
		Password:   "password456",
		GlobalRole: domain.RoleAdmin,
	}

	_, err = createUser.Execute(input2)
	if !errors.Is(err, domain.ErrUserEmailAlreadyInUse) {
		t.Errorf("expected ErrUserEmailAlreadyInUse, got %v", err)
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "invalid-email",
		Password:   "password123",
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserEmailInvalid) {
		t.Errorf("expected ErrUserEmailInvalid, got %v", err)
	}
}

func TestCreateUserInvalidRole(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Password:   "password123",
		GlobalRole: domain.UserRole("guest"),
	})
	if !errors.Is(err, domain.ErrUserRoleInvalid) {
		t.Errorf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestCreateUserPasswordTooShort(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Password:   "short",
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserPasswordTooShort) {
		t.Errorf("expected ErrUserPasswordTooShort, got %v", err)
	}
}

func TestCreateUserPropagatesUnexpectedFindByEmailError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.findByEmailErr = errors.New("database failure")
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Password:   "password123",
		GlobalRole: domain.RoleProfessor,
	})
	if err == nil || err.Error() != "database failure" {
		t.Fatalf("expected database failure, got %v", err)
	}
}

func TestCreateUserPropagatesCreateError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	mockRepo.createErr = errors.New("create failure")
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Password:   "password123",
		GlobalRole: domain.RoleProfessor,
	})
	if err == nil || err.Error() != "create failure" {
		t.Fatalf("expected create failure, got %v", err)
	}
}

func TestCreateUserRejectsRequiredPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	createUser := applicationpkg.NewCreateUser(mockRepo)

	_, err := createUser.Execute(applicationpkg.CreateUserInput{
		Name:       "John Doe",
		Email:      "john@example.com",
		Password:   "   ",
		GlobalRole: domain.RoleProfessor,
	})
	if !errors.Is(err, domain.ErrUserPasswordRequired) {
		t.Fatalf("expected ErrUserPasswordRequired, got %v", err)
	}
}
