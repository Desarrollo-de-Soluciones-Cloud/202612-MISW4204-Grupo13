package application

import (
	applicationpkg "backend/internal/auth/application"
	authDomain "backend/internal/auth/domain"
	usersDomain "backend/internal/users/domain"
	"errors"
	"testing"
)

func TestValidateTokenSuccess(t *testing.T) {
	tokenManager := NewMockTokenManager()
	tokenManager.user = &authDomain.AuthenticatedUser{
		ID:         7,
		Name:       testAuthJaneName,
		Email:      testAuthJaneEmail,
		GlobalRole: usersDomain.RoleProfessor,
	}
	validateToken := applicationpkg.NewValidateToken(tokenManager)

	output, err := validateToken.Execute(applicationpkg.ValidateTokenInput{Token: "valid-token"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.Email != testAuthJaneEmail {
		t.Fatalf("expected email %s, got %q", testAuthJaneEmail, output.Email)
	}

	if output.GlobalRole != string(usersDomain.RoleProfessor) {
		t.Fatalf("expected role professor, got %q", output.GlobalRole)
	}
}

func TestValidateTokenRequiresToken(t *testing.T) {
	tokenManager := NewMockTokenManager()
	validateToken := applicationpkg.NewValidateToken(tokenManager)

	_, err := validateToken.Execute(applicationpkg.ValidateTokenInput{})
	if !errors.Is(err, authDomain.ErrAuthTokenRequired) {
		t.Fatalf("expected ErrAuthTokenRequired, got %v", err)
	}
}

func TestValidateTokenPropagatesTokenErrors(t *testing.T) {
	tokenManager := NewMockTokenManager()
	tokenManager.parseErr = authDomain.ErrAuthTokenExpired
	validateToken := applicationpkg.NewValidateToken(tokenManager)

	_, err := validateToken.Execute(applicationpkg.ValidateTokenInput{Token: "expired-token"})
	if !errors.Is(err, authDomain.ErrAuthTokenExpired) {
		t.Fatalf("expected ErrAuthTokenExpired, got %v", err)
	}
}
