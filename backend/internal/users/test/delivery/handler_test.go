package delivery_test

import (
	authDomain "backend/internal/auth/domain"
	applicationpkg "backend/internal/users/application"
	deliverypkg "backend/internal/users/delivery"
	usersDomain "backend/internal/users/domain"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockUserRepository struct {
	users map[string]*usersDomain.User
	byID  map[uint]*usersDomain.User
	next  uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: map[string]*usersDomain.User{},
		byID:  map[uint]*usersDomain.User{},
		next:  1,
	}
}

func (m *mockUserRepository) Create(user *usersDomain.User) error {
	user.ID = m.next
	m.next++
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*usersDomain.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}

func (m *mockUserRepository) FindAll() ([]usersDomain.User, error) {
	result := make([]usersDomain.User, 0, len(m.byID))
	for _, user := range m.byID {
		result = append(result, *user)
	}
	return result, nil
}

func (m *mockUserRepository) FindAllByRole(role usersDomain.UserRole) ([]usersDomain.User, error) {
	result := make([]usersDomain.User, 0)
	for _, user := range m.byID {
		if user.GlobalRole == role {
			result = append(result, *user)
		}
	}
	return result, nil
}

func (m *mockUserRepository) FindByID(id uint) (*usersDomain.User, error) {
	if user, ok := m.byID[id]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}

func (m *mockUserRepository) Update(user *usersDomain.User) error {
	if _, ok := m.byID[user.ID]; !ok {
		return usersDomain.ErrUserNotFound
	}
	m.byID[user.ID] = user
	m.users[user.Email] = user
	return nil
}

func newUserHandlerForTest() *deliverypkg.UserHandler {
	repo := newMockUserRepository()
	return deliverypkg.NewUserHandler(
		applicationpkg.NewCreateUser(repo),
		applicationpkg.NewListUsers(repo),
		applicationpkg.NewListUsersByRole(repo),
		applicationpkg.NewGetUserByID(repo),
		applicationpkg.NewUpdateUser(repo),
		applicationpkg.NewChangeUserRole(repo),
	)
}

func TestCreateUserBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":1}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateUser(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListUsersUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users", nil)

	handler.ListUsers(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestListUsersProfessorForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users", nil)
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleProfessor,
	})

	handler.ListUsers(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestGetUserByIDBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "bad"}}

	handler.GetUserByID(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
