package delivery_test

import (
	applicationpkg "backend/internal/auth/application"
	deliverypkg "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedHelpers "backend/internal/shared/helpers"
	usersDomain "backend/internal/users/domain"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	testAuthEmail         = "john@example.com"
	testAuthSignInPath    = "/auth/sign-in"
	testAuthSecurePath    = "/secure"
	testHeaderContentType = "Content-Type"
	testApplicationJSON   = "application/json"
	errExpected200        = "expected 200, got %d"
	errExpected401        = "expected 401, got %d"
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
				Email:      testAuthEmail,
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
	req := httptest.NewRequest(http.MethodPost, testAuthSignInPath, bytes.NewBufferString(`{"email":123}`))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
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
		Email:    testAuthEmail,
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, testAuthSignInPath, bytes.NewBuffer(body))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.SignIn(c)

	if w.Code != http.StatusOK {
		t.Fatalf(errExpected200, w.Code)
	}
}

func TestGetCurrentUserUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetCurrentUser(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf(errExpected401, w.Code)
	}
}

func TestRequireAuthenticationMissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testAuthSecurePath, nil)

	handler.RequireAuthentication()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf(errExpected401, w.Code)
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

func TestSignInInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hash, _ := sharedHelpers.HashPassword("password123")
	userReader := &mockAuthUserReader{
		user: &authDomain.AuthenticatedUserCredentials{
			AuthenticatedUser: authDomain.AuthenticatedUser{
				ID:         1,
				Name:       "John",
				Email:      testAuthEmail,
				GlobalRole: usersDomain.RoleAdmin,
			},
			Password: hash,
		},
	}
	handler := deliverypkg.NewAuthHandler(
		applicationpkg.NewSignIn(userReader, &mockAuthTokenManager{}),
		applicationpkg.NewValidateToken(&mockAuthTokenManager{}),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(deliverypkg.SignInRequest{
		Email:    testAuthEmail,
		Password: "wrong-password",
	})
	req := httptest.NewRequest(http.MethodPost, testAuthSignInPath, bytes.NewBuffer(body))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.SignIn(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf(errExpected401, w.Code)
	}
}

func TestSignInInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewAuthHandler(
		applicationpkg.NewSignIn(&mockAuthUserReader{err: errors.New("boom")}, &mockAuthTokenManager{}),
		applicationpkg.NewValidateToken(&mockAuthTokenManager{}),
	)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(deliverypkg.SignInRequest{
		Email:    testAuthEmail,
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, testAuthSignInPath, bytes.NewBuffer(body))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.SignIn(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetCurrentUserSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		Name:       "John",
		Email:      testAuthEmail,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.GetCurrentUser(c)

	if w.Code != http.StatusOK {
		t.Fatalf(errExpected200, w.Code)
	}
}

func TestRequireAuthenticationInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewAuthHandler(
		applicationpkg.NewSignIn(&mockAuthUserReader{}, &mockAuthTokenManager{}),
		applicationpkg.NewValidateToken(&mockAuthTokenManager{err: authDomain.ErrAuthTokenInvalid}),
	)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testAuthSecurePath, nil)
	c.Request.Header.Set("Authorization", "Bearer bad-token")

	handler.RequireAuthentication()(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf(errExpected401, w.Code)
	}
}

func TestRequireAuthenticationSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	user := &authDomain.AuthenticatedUser{
		ID:         1,
		Name:       "John",
		Email:      testAuthEmail,
		GlobalRole: usersDomain.RoleAdmin,
	}
	handler := deliverypkg.NewAuthHandler(
		applicationpkg.NewSignIn(&mockAuthUserReader{}, &mockAuthTokenManager{}),
		applicationpkg.NewValidateToken(&mockAuthTokenManager{user: user}),
	)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testAuthSecurePath, nil)
	c.Request.Header.Set("Authorization", "Bearer ok-token")

	handler.RequireAuthentication()(c)

	if w.Code != http.StatusOK {
		t.Fatalf(errExpected200, w.Code)
	}
	if _, ok := deliverypkg.GetCurrentUser(c); !ok {
		t.Fatal("expected current user in context")
	}
}

func TestRequireRolesUnauthorizedWithoutUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.RequireRoles(usersDomain.RoleAdmin)(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf(errExpected401, w.Code)
	}
}

func TestRequireRolesSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newAuthHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.RequireRoles(usersDomain.RoleAdmin)(c)

	if w.Code != http.StatusOK {
		t.Fatalf(errExpected200, w.Code)
	}
}
