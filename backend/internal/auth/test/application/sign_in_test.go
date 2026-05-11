package application

import (
	applicationpkg "backend/internal/auth/application"
	authDomain "backend/internal/auth/domain"
	sharedHelpers "backend/internal/shared/helpers"
	usersDomain "backend/internal/users/domain"
	"errors"
	"testing"
)

const (
	testAuthEmailJohn   = "john@example.com"
	testAuthEmailJohnUC = "JOHN@example.com"
	testAuthUserJohn    = "John Doe"
	testAuthPassword123 = "password123"
	testAuthMissingEmail = "missing@example.com"
	errExpectedHash     = "expected password hash, got %v"
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

	passwordHash, err := sharedHelpers.HashPassword(testAuthPassword123)
	if err != nil {
		t.Fatalf(errExpectedHash, err)
	}

	userReader.users[testAuthEmailJohn] = &authDomain.AuthenticatedUserCredentials{
		AuthenticatedUser: authDomain.AuthenticatedUser{
			ID:         1,
			Name:       testAuthUserJohn,
			Email:      testAuthEmailJohn,
			GlobalRole: usersDomain.RoleProfessor,
		},
		Password: passwordHash,
	}

	output, err := signIn.Execute(applicationpkg.SignInInput{
		Email:    testAuthEmailJohnUC,
		Password: testAuthPassword123,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.AccessToken != "mock-token" {
		t.Fatalf("expected mock token, got %q", output.AccessToken)
	}

	if output.User.Email != testAuthEmailJohn {
		t.Fatalf("expected normalized email, got %q", output.User.Email)
	}
}

func TestSignInInvalidCredentialsWhenUserDoesNotExist(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	_, err := signIn.Execute(applicationpkg.SignInInput{
		Email:    testAuthMissingEmail,
		Password: testAuthPassword123,
	})
	if !errors.Is(err, authDomain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestSignInInvalidCredentialsWhenPasswordDoesNotMatch(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	passwordHash, err := sharedHelpers.HashPassword(testAuthPassword123)
	if err != nil {
		t.Fatalf(errExpectedHash, err)
	}

	userReader.users[testAuthEmailJohn] = &authDomain.AuthenticatedUserCredentials{
		AuthenticatedUser: authDomain.AuthenticatedUser{
			ID:         1,
			Name:       testAuthUserJohn,
			Email:      testAuthEmailJohn,
			GlobalRole: usersDomain.RoleAdmin,
		},
		Password: passwordHash,
	}

	_, err = signIn.Execute(applicationpkg.SignInInput{
		Email:    testAuthEmailJohn,
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
		Password: testAuthPassword123,
	})
	if !errors.Is(err, authDomain.ErrAuthEmailInvalid) {
		t.Fatalf("expected ErrAuthEmailInvalid, got %v", err)
	}
}

func TestSignInRejectsShortPassword(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	_, err := signIn.Execute(applicationpkg.SignInInput{
		Email:    testAuthEmailJohn,
		Password: "short",
	})
	if !errors.Is(err, authDomain.ErrAuthPasswordTooShort) {
		t.Fatalf("expected ErrAuthPasswordTooShort, got %v", err)
	}
}

func TestSignInPropagatesUnexpectedReaderError(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	userReader.users = nil
	readerErr := errors.New("reader failure")
	userReaderWithErr := &MockUserReaderWithError{err: readerErr}
	signIn = applicationpkg.NewSignIn(userReaderWithErr, tokenManager)

	_, err := signIn.Execute(applicationpkg.SignInInput{
		Email:    testAuthEmailJohn,
		Password: testAuthPassword123,
	})
	if !errors.Is(err, readerErr) {
		t.Fatalf("expected reader error, got %v", err)
	}
}

func TestSignInPropagatesTokenGenerationError(t *testing.T) {
	userReader := NewMockUserReader()
	tokenManager := NewMockTokenManager()
	tokenManager.generateErr = errors.New("token generation failed")
	signIn := applicationpkg.NewSignIn(userReader, tokenManager)

	passwordHash, err := sharedHelpers.HashPassword(testAuthPassword123)
	if err != nil {
		t.Fatalf(errExpectedHash, err)
	}

	userReader.users[testAuthEmailJohn] = &authDomain.AuthenticatedUserCredentials{
		AuthenticatedUser: authDomain.AuthenticatedUser{
			ID:         1,
			Name:       testAuthUserJohn,
			Email:      testAuthEmailJohn,
			GlobalRole: usersDomain.RoleProfessor,
		},
		Password: passwordHash,
	}

	_, err = signIn.Execute(applicationpkg.SignInInput{
		Email:    testAuthEmailJohn,
		Password: testAuthPassword123,
	})
	if err == nil || err.Error() != "token generation failed" {
		t.Fatalf("expected token generation error, got %v", err)
	}
}

type MockUserReaderWithError struct {
	err error
}

func (m *MockUserReaderWithError) GetByEmail(email string) (*authDomain.AuthenticatedUserCredentials, error) {
	return nil, m.err
}
