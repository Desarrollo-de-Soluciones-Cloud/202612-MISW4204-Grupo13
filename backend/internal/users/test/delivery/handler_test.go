package delivery

import (
	authDomain "backend/internal/auth/domain"
	usersApplication "backend/internal/users/application"
	usersDelivery "backend/internal/users/delivery"
	usersDomain "backend/internal/users/domain"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockUserRepository struct {
	users            map[string]*usersDomain.User
	byID             map[uint]*usersDomain.User
	nextID           uint
	findByIDErr      error
	findAllErr       error
	findAllByRoleErr error
	updateErr        error
	createErr        error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[string]*usersDomain.User),
		byID:   make(map[uint]*usersDomain.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) Create(user *usersDomain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepository) FindByID(id uint) (*usersDomain.User, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	if user, ok := m.byID[id]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}

func (m *mockUserRepository) FindByEmail(email string) (*usersDomain.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}

func (m *mockUserRepository) FindAll() ([]usersDomain.User, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	users := make([]usersDomain.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *mockUserRepository) FindAllByRole(role usersDomain.UserRole) ([]usersDomain.User, error) {
	if m.findAllByRoleErr != nil {
		return nil, m.findAllByRoleErr
	}
	users := make([]usersDomain.User, 0)
	for _, user := range m.users {
		if user.GlobalRole == role {
			users = append(users, *user)
		}
	}
	return users, nil
}

func (m *mockUserRepository) Update(user *usersDomain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.byID[user.ID]; !ok {
		return usersDomain.ErrUserNotFound
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

func newUserHandler(repo *mockUserRepository) *usersDelivery.UserHandler {
	return usersDelivery.NewUserHandler(
		usersApplication.NewCreateUser(repo),
		usersApplication.NewListUsers(repo),
		usersApplication.NewListUsersByRole(repo),
		usersApplication.NewGetUserByID(repo),
		usersApplication.NewUpdateUser(repo),
		usersApplication.NewChangeUserRole(repo),
	)
}

func seedUser(t *testing.T, repo *mockUserRepository, name, email string, role usersDomain.UserRole) uint {
	t.Helper()

	output, err := usersApplication.NewCreateUser(repo).Execute(usersApplication.CreateUserInput{
		Name:       name,
		Email:      email,
		Password:   "password123",
		GlobalRole: role,
	})
	if err != nil {
		t.Fatalf("expected seed user, got %v", err)
	}
	return output.ID
}

func TestCreateUserBindingError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"name":"","email":"bad-email","password":"123","global_role":""}`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestCreateUserSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"name":"John Doe","email":"john@example.com","password":"password123","global_role":"professor"}`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}
}

func TestCreateUserInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	repo.createErr = errors.New("db error")
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"name":"John Doe","email":"john@example.com","password":"password123","global_role":"professor"}`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestCreateUserConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"name":"Jane Doe","email":"john@example.com","password":"password123","global_role":"assistant"}`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", recorder.Code)
	}
}

func TestCreateUserDomainValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"name":"John Doe","email":"john@example.com","password":"password123","global_role":"guest"}`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestListUsersWithoutCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	context.Request = request

	handler.ListUsers(context)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestListUsersProfessorMustFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	context.Request = request
	context.Set("current_user", authDomain.AuthenticatedUser{GlobalRole: usersDomain.RoleProfessor})

	handler.ListUsers(context)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestListUsersProfessorAllowedFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "Monitor One", "monitor@example.com", usersDomain.RoleMonitor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/users?role=monitor", nil)
	context.Request = request
	context.Set("current_user", authDomain.AuthenticatedUser{GlobalRole: usersDomain.RoleProfessor})

	handler.ListUsers(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestListUsersProfessorRejectedFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/users?role=admin", nil)
	context.Request = request
	context.Set("current_user", authDomain.AuthenticatedUser{GlobalRole: usersDomain.RoleProfessor})

	handler.ListUsers(context)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestListUsersAdminSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	context.Request = request
	context.Set("current_user", authDomain.AuthenticatedUser{GlobalRole: usersDomain.RoleAdmin})

	handler.ListUsers(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestListUsersInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	repo.findAllErr = errors.New("db error")
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	context.Request = request
	context.Set("current_user", authDomain.AuthenticatedUser{GlobalRole: usersDomain.RoleAdmin})

	handler.ListUsers(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestListUsersAdminWithRoleFilterDelegatesSuccessfully(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "Monitor One", "monitor@example.com", usersDomain.RoleMonitor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/users?role=monitor", nil)
	context.Request = request
	context.Set("current_user", authDomain.AuthenticatedUser{GlobalRole: usersDomain.RoleAdmin})

	handler.ListUsers(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestListUsersByRoleInvalidRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	handler.ListUsersByRole(context, "guest")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestListUsersByRoleInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	repo.findAllByRoleErr = errors.New("db error")
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	handler.ListUsersByRole(context, "monitor")

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestGetUserByIDInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "bad"}}

	handler.GetUserByID(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestGetUserByIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.GetUserByID(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "99"}}

	handler.GetUserByID(context)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestGetUserByIDInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	repo.findByIDErr = errors.New("db error")
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.GetUserByID(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestCreateUserBindingErrorRequiredAndMax(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	longPassword := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstu"
	requestBody := bytes.NewBufferString(`{"name":"","email":"","password":"` + longPassword + `","global_role":""}`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestCreateUserBindingErrorWithMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"name":`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestCreateUserBindingErrorNameTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	longName := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	requestBody := bytes.NewBufferString(`{"name":"` + longName + `","email":"john@example.com","password":"password123","global_role":"professor"}`)
	request, _ := http.NewRequest(http.MethodPost, "/users", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.CreateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestUpdateUserBindingError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	requestBody := bytes.NewBufferString(`{"name":"","email":"bad-email","global_role":""}`)
	request, _ := http.NewRequest(http.MethodPut, "/users/1", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.UpdateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestUpdateUserInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "bad"}}

	handler.UpdateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestUpdateUserConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	seedUser(t, repo, "Jane Doe", "jane@example.com", usersDomain.RoleMonitor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "2"}}
	requestBody := bytes.NewBufferString(`{"name":"Jane Doe","email":"john@example.com","global_role":"monitor"}`)
	request, _ := http.NewRequest(http.MethodPut, "/users/2", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.UpdateUser(context)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", recorder.Code)
	}
}

func TestUpdateUserValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	requestBody := bytes.NewBufferString(`{"name":"John Doe","email":"john@example.com","global_role":"guest"}`)
	request, _ := http.NewRequest(http.MethodPut, "/users/1", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.UpdateUser(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestUpdateUserInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	repo.updateErr = errors.New("db error")
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	requestBody := bytes.NewBufferString(`{"name":"John Doe","email":"john@example.com","global_role":"professor"}`)
	request, _ := http.NewRequest(http.MethodPut, "/users/1", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.UpdateUser(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	requestBody := bytes.NewBufferString(`{"name":"Jane Doe","email":"jane@example.com","global_role":"assistant"}`)
	request, _ := http.NewRequest(http.MethodPut, "/users/1", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.UpdateUser(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "99"}}
	requestBody := bytes.NewBufferString(`{"name":"Jane Doe","email":"jane@example.com","global_role":"assistant"}`)
	request, _ := http.NewRequest(http.MethodPut, "/users/99", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.UpdateUser(context)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestChangeUserRoleBindingError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	requestBody := bytes.NewBufferString(`{"global_role":""}`)
	request, _ := http.NewRequest(http.MethodPatch, "/users/1/role", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.ChangeUserRole(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestChangeUserRoleInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "bad"}}

	handler.ChangeUserRole(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestChangeUserRoleValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	requestBody := bytes.NewBufferString(`{"global_role":"guest"}`)
	request, _ := http.NewRequest(http.MethodPatch, "/users/1/role", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.ChangeUserRole(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestChangeUserRoleNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := newUserHandler(newMockUserRepository())
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "99"}}
	requestBody := bytes.NewBufferString(`{"global_role":"assistant"}`)
	request, _ := http.NewRequest(http.MethodPatch, "/users/99/role", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.ChangeUserRole(context)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestChangeUserRoleSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newMockUserRepository()
	seedUser(t, repo, "John Doe", "john@example.com", usersDomain.RoleProfessor)
	handler := newUserHandler(repo)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Params = gin.Params{{Key: "id", Value: "1"}}
	requestBody := bytes.NewBufferString(`{"global_role":"assistant"}`)
	request, _ := http.NewRequest(http.MethodPatch, "/users/1/role", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.ChangeUserRole(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}
