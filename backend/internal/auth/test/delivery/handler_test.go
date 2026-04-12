package delivery

import (
	applicationpkg "backend/internal/auth/application"
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedHelpers "backend/internal/shared/helpers"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	usersDomain "backend/internal/users/domain"

	"github.com/gin-gonic/gin"
)

type mockUserReader struct {
	users map[string]*authDomain.AuthenticatedUserCredentials
	err   error
}

func (m *mockUserReader) GetByEmail(email string) (*authDomain.AuthenticatedUserCredentials, error) {
	if m.err != nil {
		return nil, m.err
	}
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, usersDomain.ErrUserNotFound
}

type mockTokenManager struct {
	token       *authDomain.AuthToken
	user        *authDomain.AuthenticatedUser
	generateErr error
	parseErr    error
}

func (m *mockTokenManager) GenerateToken(user *authDomain.AuthenticatedUser) (*authDomain.AuthToken, error) {
	if m.generateErr != nil {
		return nil, m.generateErr
	}
	m.user = user
	return m.token, nil
}

func (m *mockTokenManager) ParseToken(token string) (*authDomain.AuthenticatedUser, error) {
	if m.parseErr != nil {
		return nil, m.parseErr
	}
	return m.user, nil
}

func newAuthHandler(t *testing.T) (*authDelivery.AuthHandler, *mockUserReader, *mockTokenManager) {
	t.Helper()

	passwordHash, err := sharedHelpers.HashPassword("password123")
	if err != nil {
		t.Fatalf("expected password hash, got %v", err)
	}

	userReader := &mockUserReader{
		users: map[string]*authDomain.AuthenticatedUserCredentials{
			"john@example.com": {
				AuthenticatedUser: authDomain.AuthenticatedUser{
					ID:         1,
					Name:       "John Doe",
					Email:      "john@example.com",
					GlobalRole: usersDomain.RoleProfessor,
				},
				Password: passwordHash,
			},
		},
	}
	tokenManager := &mockTokenManager{
		token: &authDomain.AuthToken{
			AccessToken: "valid-token",
			TokenType:   authDomain.TokenTypeBearer,
			ExpiresIn:   3600,
		},
		user: &authDomain.AuthenticatedUser{
			ID:         1,
			Name:       "John Doe",
			Email:      "john@example.com",
			GlobalRole: usersDomain.RoleProfessor,
		},
	}

	return authDelivery.NewAuthHandler(
		applicationpkg.NewSignIn(userReader, tokenManager),
		applicationpkg.NewValidateToken(tokenManager),
	), userReader, tokenManager
}

func TestSignInSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, _ := newAuthHandler(t)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"email":"john@example.com","password":"password123"}`)
	request, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.SignIn(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestSignInBindingError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, _ := newAuthHandler(t)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"email":"bad-email","password":"123"}`)
	request, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.SignIn(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var payload map[string][]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid error response, got %v", err)
	}
	if len(payload["errors"]) != 2 {
		t.Fatalf("expected 2 binding errors, got %v", payload["errors"])
	}
}

func TestSignInRejectsInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, userReader, _ := newAuthHandler(t)
	userReader.users = map[string]*authDomain.AuthenticatedUserCredentials{}

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"email":"missing@example.com","password":"password123"}`)
	request, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.SignIn(context)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestSignInReturnsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, userReader, _ := newAuthHandler(t)
	userReader.err = errors.New("db failure")

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	requestBody := bytes.NewBufferString(`{"email":"john@example.com","password":"password123"}`)
	request, _ := http.NewRequest(http.MethodPost, "/auth/sign-in", requestBody)
	request.Header.Set("Content-Type", "application/json")
	context.Request = request

	handler.SignIn(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestGetCurrentUserWithoutUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, _ := newAuthHandler(t)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	handler.GetCurrentUser(context)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestRequireAuthenticationRejectsExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, tokenManager := newAuthHandler(t)
	tokenManager.parseErr = authDomain.ErrAuthTokenExpired

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/secure", nil)
	request.Header.Set("Authorization", "Bearer expired-token")
	context.Request = request

	handler.RequireAuthentication()(context)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestRequireAuthenticationUnexpectedError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, tokenManager := newAuthHandler(t)
	tokenManager.parseErr = errors.New("backend failure")

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	request, _ := http.NewRequest(http.MethodGet, "/secure", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	context.Request = request

	handler.RequireAuthentication()(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
}

func TestRequireRolesRejectsMissingCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler, _, _ := newAuthHandler(t)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	handler.RequireRoles(usersDomain.RoleAdmin)(context)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
