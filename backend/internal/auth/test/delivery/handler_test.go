package delivery_test

import (
	applicationpkg "backend/internal/auth/application"
	deliverypkg "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedHelpers "backend/internal/shared/helpers"
	usersDomain "backend/internal/users/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockAuthUserReader struct {
	user *authDomain.AuthenticatedUserCredentials
	err  error
}

func (m *mockAuthUserReader) GetByEmail(email string) (*authDomain.AuthenticatedUserCredentials, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

type mockAuthTokenManager struct {
	token *authDomain.AuthToken
	user  *authDomain.AuthenticatedUser
	err   error
}

func (m *mockAuthTokenManager) GenerateToken(user *authDomain.AuthenticatedUser) (*authDomain.AuthToken, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.user = user
	if m.token == nil {
		m.token = &authDomain.AuthToken{
			AccessToken: "mock-token",
			TokenType:   authDomain.TokenTypeBearer,
			ExpiresIn:   3600,
		}
	}
	return m.token, nil
}

func (m *mockAuthTokenManager) ParseToken(token string) (*authDomain.AuthenticatedUser, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func newAuthHandlerForTest(t *testing.T) *deliverypkg.AuthHandler {
	t.Helper()

	hash, err := sharedHelpers.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}

	userReader := &mockAuthUserReader{
		user: &authDomain.AuthenticatedUserCredentials{
			AuthenticatedUser: authDomain.AuthenticatedUser{
				ID:         1,
				Name:       "John",
				Email:      "john@example.com",
				GlobalRole: usersDomain.RoleAdmin,
			},
			Password: hash,
		},
	}
	tokenManager := &mockAuthTokenManager{}

	return deliverypkg.NewAuthHandler(
		applicationpkg.NewSignIn(userReader, tokenManager),
		applicationpkg.NewValidateToken(tokenManager),
	)
}

func TestSignInBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(`{"email":123}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.SignIn(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSignInSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(deliverypkg.SignInRequest{
		Email:    "john@example.com",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.SignIn(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetCurrentUserUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetCurrentUser(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuthenticationMissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/secure", nil)

	handler.RequireAuthentication()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireRolesForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleProfessor,
	})

	handler.RequireRoles(usersDomain.RoleAdmin)(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}
