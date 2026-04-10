package application

import (
	"backend/internal/users/application"
	"backend/internal/users/domain"
	"testing"
)

type MockListUsersRepository struct {
    users []domain.User
}

func (m *MockListUsersRepository) Create(user *domain.User) error {
    user.ID = uint(len(m.users) + 1)
    m.users = append(m.users, *user)
    return nil
}

func (m *MockListUsersRepository) FindByEmail(email string) (*domain.User, error) {
    for _, u := range m.users {
        if u.Email == email {
            return &u, nil
        }
    }
    return nil, domain.ErrUserNotFound
}

func (m *MockListUsersRepository) FindAll() ([]domain.User, error) {
    return m.users, nil
}

func (m *MockListUsersRepository) FindByID(id uint) (*domain.User, error) {
    for i := range m.users {
        if m.users[i].ID == id {
            return &m.users[i], nil
        }
    }
    return nil, domain.ErrUserNotFound
}

func TestListUsersEmpty(t *testing.T) {
    mockRepo := &MockListUsersRepository{users: []domain.User{}}
    listUsers := application.NewListUsers(mockRepo)
    output, err := listUsers.Execute()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(output.Users) != 0 {
        t.Errorf("Expected 0 users, got %d", len(output.Users))
    }
}

func TestListUsersWithData(t *testing.T) {
    mockRepo := &MockListUsersRepository{
        users: []domain.User{
            {
                ID:        1,
                Name:      "John Doe",
                Email:     "john@example.com",
                Password:  "password",
            },
            {
                ID:        2,
                Name:      "Jane Doe",
                Email:     "jane@example.com",
                Password:  "password",
            },
        },
    }
    listUsers := application.NewListUsers(mockRepo)
    output, err := listUsers.Execute()
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(output.Users) != 2 {
        t.Errorf("Expected 2 users, got %d", len(output.Users))
    }
    if output.Users[0].Name != "John Doe" {
        t.Errorf("Expected first user name 'John Doe', got '%s'", output.Users[0].Name)
    }
    if output.Users[1].Name != "Jane Doe" {
        t.Errorf("Expected second user name 'Jane Doe', got '%s'", output.Users[1].Name)
    }
}