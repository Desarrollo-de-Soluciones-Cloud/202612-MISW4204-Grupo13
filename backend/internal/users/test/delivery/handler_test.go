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

const (
	testUsersPath         = "/users"
	testHeaderContentType = "Content-Type"
	testApplicationJSON   = "application/json"
	errExpected200        = "expected 200, got %d"
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
	req := httptest.NewRequest(http.MethodPost, testUsersPath, bytes.NewBufferString(`{"name":1}`))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
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
	c.Request = httptest.NewRequest(http.MethodGet, testUsersPath, nil)

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
	c.Request = httptest.NewRequest(http.MethodGet, testUsersPath, nil)
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

func TestCreateUserConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()

	firstBody := bytes.NewBufferString(`{"name":"Ana Gomez","email":"ana@example.com","password":"Password123","global_role":"professor"}`)
	firstReq := httptest.NewRequest(http.MethodPost, testUsersPath, firstBody)
	firstReq.Header.Set(testHeaderContentType, testApplicationJSON)
	firstW := httptest.NewRecorder()
	firstC, _ := gin.CreateTestContext(firstW)
	firstC.Request = firstReq
	handler.CreateUser(firstC)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, testUsersPath, bytes.NewBufferString(`{"name":"Ana Gomez","email":"ana@example.com","password":"Password123","global_role":"professor"}`))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateUser(c)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestListUsersAdminSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()
	repo := newMockUserRepository()
	_ = repo.Create(&usersDomain.User{ID: 1, Name: "Ana", Email: "ana@example.com", GlobalRole: usersDomain.RoleProfessor})
	_ = repo.Create(&usersDomain.User{ID: 2, Name: "Luis", Email: "luis@example.com", GlobalRole: usersDomain.RoleMonitor})
	handler = deliverypkg.NewUserHandler(
		applicationpkg.NewCreateUser(repo),
		applicationpkg.NewListUsers(repo),
		applicationpkg.NewListUsersByRole(repo),
		applicationpkg.NewGetUserByID(repo),
		applicationpkg.NewUpdateUser(repo),
		applicationpkg.NewChangeUserRole(repo),
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testUsersPath, nil)
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})

	handler.ListUsers(c)

	if w.Code != http.StatusOK {
		t.Fatalf(errExpected200, w.Code)
	}
}

func TestListUsersProfessorRoleFilterAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()
	repo := newMockUserRepository()
	_ = repo.Create(&usersDomain.User{Name: "Luis", Email: "luis@example.com", GlobalRole: usersDomain.RoleMonitor})
	handler = deliverypkg.NewUserHandler(
		applicationpkg.NewCreateUser(repo),
		applicationpkg.NewListUsers(repo),
		applicationpkg.NewListUsersByRole(repo),
		applicationpkg.NewGetUserByID(repo),
		applicationpkg.NewUpdateUser(repo),
		applicationpkg.NewChangeUserRole(repo),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, testUsersPath+"?role=monitor", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleProfessor})

	handler.ListUsers(c)

	if w.Code != http.StatusOK {
		t.Fatalf(errExpected200, w.Code)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newUserHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users/99", nil)
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	handler.GetUserByID(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestChangeUserRoleSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockUserRepository()
	user, _ := usersDomain.NewUser("Ana Gomez", "ana@example.com", "Password123", usersDomain.RoleProfessor)
	_ = repo.Create(user)
	handler := deliverypkg.NewUserHandler(
		applicationpkg.NewCreateUser(repo),
		applicationpkg.NewListUsers(repo),
		applicationpkg.NewListUsersByRole(repo),
		applicationpkg.NewGetUserByID(repo),
		applicationpkg.NewUpdateUser(repo),
		applicationpkg.NewChangeUserRole(repo),
	)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/users/1/role", bytes.NewBufferString(`{"global_role":"assistant"}`))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.ChangeUserRole(c)

	if w.Code != http.StatusOK {
		t.Fatalf(errExpected200, w.Code)
	}
}
