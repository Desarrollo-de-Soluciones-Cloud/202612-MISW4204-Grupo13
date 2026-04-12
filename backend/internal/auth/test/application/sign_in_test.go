package application

import (
	applicationpkg "backend/internal/auth/application"
	authDomain "backend/internal/auth/domain"
	sharedHelpers "backend/internal/shared/helpers"
	usersDomain "backend/internal/users/domain"
	"errors"
	"testing"
)

type MockUserReader struct {
	users map[string]*authDomain.AuthenticatedUserCredentials
}

func NewMockUserReader() *MockUserReader {
	return &MockUserReader{
		users: make(map[string]*authDomain.AuthenticatedUserCredentials),
	}
}

func (m *MockUserReader) GetByEmail(email string) (*authDomain.AuthenticatedUserCredentials, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}

	return nil, usersDomain.ErrUserNotFound
}

type MockTokenManager struct {
	token       *authDomain.AuthToken
	user        *authDomain.AuthenticatedUser
	generateErr error
	parseErr    error
}

func NewMockTokenManager() *MockTokenManager {
	return &MockTokenManager{
		token: &authDomain.AuthToken{
			AccessToken: "mock-token",
			TokenType:   authDomain.TokenTypeBearer,
			ExpiresIn:   3600,
		},
	}
}

func (m *MockTokenManager) GenerateToken(user *authDomain.AuthenticatedUser) (*authDomain.AuthToken, error) {
	if m.generateErr != nil {
		return nil, m.generateErr
	}

	m.user = user
	return m.token, nil
}

func (m *MockTokenManager) ParseToken(token string) (*authDomain.AuthenticatedUser, error) {
	if m.parseErr != nil {
		return nil, m.parseErr
	}

	return m.user, nil
}

func TestSignInSuccess(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	passwordHash, err := sharedHelpers.HashPassword("password123")
	if err != nil {
		t.Fatalf("expected password hash, got %v", err)
	}

	userReader.users["john@example.com"] = &authDomain.AuthenticatedUserCredentials{
		AuthenticatedUser: authDomain.AuthenticatedUser{
			ID:         1,
			Name:       "John Doe",
			Email:      "john@example.com",
			GlobalRole: usersDomain.RoleProfessor,
		},
		Password: passwordHash,
	}

	output, err := signIn.Execute(applicationpkg.SignInInput{
		Email:    "JOHN@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.AccessToken != "mock-token" {
		t.Fatalf("expected mock token, got %q", output.AccessToken)
	}

	if output.User.Email != "john@example.com" {
		t.Fatalf("expected normalized email, got %q", output.User.Email)
	}
}

func TestSignInInvalidCredentialsWhenUserDoesNotExist(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	_, err := signIn.Execute(applicationpkg.SignInInput{
		Email:    "missing@example.com",
		Password: "password123",
	})
	if !errors.Is(err, authDomain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestSignInInvalidCredentialsWhenPasswordDoesNotMatch(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	passwordHash, err := sharedHelpers.HashPassword("password123")
	if err != nil {
		t.Fatalf("expected password hash, got %v", err)
	}

	userReader.users["john@example.com"] = &authDomain.AuthenticatedUserCredentials{
		AuthenticatedUser: authDomain.AuthenticatedUser{
			ID:         1,
			Name:       "John Doe",
			Email:      "john@example.com",
			GlobalRole: usersDomain.RoleAdmin,
		},
		Password: passwordHash,
	}

	_, err = signIn.Execute(applicationpkg.SignInInput{
		Email:    "john@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, authDomain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestSignInRejectsInvalidEmail(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	_, err := signIn.Execute(applicationpkg.SignInInput{
		Email:    "invalid-email",
		Password: "password123",
	})
	if !errors.Is(err, authDomain.ErrAuthEmailInvalid) {
		t.Fatalf("expected ErrAuthEmailInvalid, got %v", err)
	}
}
