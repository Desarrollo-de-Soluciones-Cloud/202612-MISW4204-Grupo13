package application

import (
	"backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"testing"
)

type MockUserRepository struct {
    users       map[string]*domain.User
    nextID      uint
    findByEmailFn func(email string) (*domain.User, error)
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users:  make(map[string]*domain.User),
        nextID: 1,
    }
}

func (m *MockUserRepository) Create(user *domain.User) error {
    user.ID = m.nextID
    m.nextID++
    m.users[user.Email] = user
    return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
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

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
    for _, u := range m.users {
        if u.ID == id {
            return u, nil
        }
    }
    return nil, domain.ErrUserNotFound
}

func TestCreateUserSuccess(t *testing.T) {
    mockRepo := NewMockUserRepository()
    createUser := application.NewCreateUser(mockRepo)
    input := application.CreateUserInput{
        Name:     "John Doe",
        Email:    "john@example.com",
        Password: "password123",
    }
    output, err := createUser.Execute(input)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output.Name != "John Doe" {
        t.Errorf("Expected name 'John Doe', got '%s'", output.Name)
    }
    if output.Email != "john@example.com" {
        t.Errorf("Expected email 'john@example.com', got '%s'", output.Email)
    }
}

func TestCreateUserEmptyName(t *testing.T) {
    mockRepo := NewMockUserRepository()
    createUser := application.NewCreateUser(mockRepo)
    input := application.CreateUserInput{
        Name:     "",
        Email:    "john@example.com",
        Password: "password123",
    }
    _, err := createUser.Execute(input)
    if !errors.Is(err, domain.ErrInvalidInput) {
        t.Errorf("Expected ErrInvalidInput, got %v", err)
    }
}

func TestCreateUserDuplicateEmail(t *testing.T) {
    mockRepo := NewMockUserRepository()
    createUser := application.NewCreateUser(mockRepo)
    input := application.CreateUserInput{
        Name:     "John Doe",
        Email:    "john@example.com",
        Password: "password123",
    }
    _, err := createUser.Execute(input)
    if err != nil {
        t.Fatalf("First creation should succeed, got %v", err)
    }
    input2 := application.CreateUserInput{
        Name:     "Jane Doe",
        Email:    "john@example.com",
        Password: "password456",
    }
    _, err = createUser.Execute(input2)
    if !errors.Is(err, domain.ErrUserAlreadyExists) {
        t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
    }
}