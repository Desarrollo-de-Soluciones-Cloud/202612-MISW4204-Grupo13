package application

import (
	applicationpkg "backend/internal/users/application"
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
	for i := range m.users {
		if m.users[i].Email == email {
			return &m.users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockListUsersRepository) FindAll() ([]domain.User, error) {
	return m.users, nil
}

func (m *MockListUsersRepository) FindAllByRole(role domain.UserRole) ([]domain.User, error) {
	users := make([]domain.User, 0)
	for _, user := range m.users {
		if user.GlobalRole == role {
			users = append(users, user)
		}
	}
	return users, nil
}

func (m *MockListUsersRepository) FindByID(id uint) (*domain.User, error) {
	for i := range m.users {
		if m.users[i].ID == id {
			return &m.users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockListUsersRepository) Update(user *domain.User) error {
	for i := range m.users {
		if m.users[i].ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return domain.ErrUserNotFound
}

func TestListUsersEmpty(t *testing.T) {
	mockRepo := &MockListUsersRepository{users: []domain.User{}}
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
	mockRepo := &MockListUsersRepository{
		users: []domain.User{
			{
				ID:         1,
				Name:       "John Doe",
				Email:      "john@example.com",
				Password:   "password-hash",
				GlobalRole: domain.RoleProfessor,
			},
			{
				ID:         2,
				Name:       "Jane Doe",
				Email:      "jane@example.com",
				Password:   "password-hash",
				GlobalRole: domain.RoleAdmin,
			},
		},
	}
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
