package infrastructure

import (
	authInfrastructure "backend/internal/auth/infrastructure"
	usersApplication "backend/internal/users/application"
	usersDomain "backend/internal/users/domain"
	"errors"
	"testing"
)

const testReaderEmailJohn = "john@example.com"

type mockUserRepository struct {
	users map[string]*usersDomain.User
}

func (m *mockUserRepository) Create(user *usersDomain.User) error { return nil }
func (m *mockUserRepository) FindByID(id uint) (*usersDomain.User, error) {
	return nil, usersDomain.ErrUserNotFound
}
func (m *mockUserRepository) FindByEmail(email string) (*usersDomain.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}
func (m *mockUserRepository) FindAll() ([]usersDomain.User, error) { return nil, nil }
func (m *mockUserRepository) FindAllByRole(role usersDomain.UserRole) ([]usersDomain.User, error) {
	return nil, nil
}
func (m *mockUserRepository) Update(user *usersDomain.User) error { return nil }

func TestGetByEmailSuccess(t *testing.T) {
	repo := &mockUserRepository{
		users: map[string]*usersDomain.User{
			testReaderEmailJohn: {
				ID:         1,
				Name:       "John Doe",
				Email:      testReaderEmailJohn,
				Password:   "hashed-password",
				GlobalRole: usersDomain.RoleProfessor,
			},
		},
	}

	getUserByEmail := usersApplication.NewGetUserByEmail(repo)
	reader := authInfrastructure.NewUserReader(getUserByEmail)
	user, err := reader.GetByEmail(testReaderEmailJohn)
	if err != nil {
		t.Fatalf("expected user, got %v", err)
	}
	if user.Password != "hashed-password" {
		t.Fatalf("expected hashed password, got %q", user.Password)
	}
}

func TestGetByEmailPropagatesError(t *testing.T) {
	repo := &mockUserRepository{users: map[string]*usersDomain.User{}}
	getUserByEmail := usersApplication.NewGetUserByEmail(repo)
	reader := authInfrastructure.NewUserReader(getUserByEmail)

	_, err := reader.GetByEmail("missing@example.com")
	if !errors.Is(err, usersDomain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
