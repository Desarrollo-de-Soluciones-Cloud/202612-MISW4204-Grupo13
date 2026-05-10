package domain

import (
	domainpkg "backend/internal/users/domain"
	"errors"
	"strings"
	"testing"
)

func TestNewUserSuccess(t *testing.T) {
	user, err := domainpkg.NewUser(" "+testDomainUserJohnName+" ", " John@Example.com ", testDomainUserPassword, domainpkg.RoleProfessor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != testDomainUserJohnName {
		t.Fatalf("expected normalized name, got %q", user.Name)
	}
	if user.Email != testDomainUserJohnEmail {
		t.Fatalf("expected normalized email, got %q", user.Email)
	}
	if user.GlobalRole != domainpkg.RoleProfessor {
		t.Fatalf("expected role %q, got %q", domainpkg.RoleProfessor, user.GlobalRole)
	}
	if user.Password != testDomainUserPassword {
		t.Fatalf("expected password to remain unchanged in domain constructor, got %q", user.Password)
	}
}

func TestNewUserRejectsInvalidRole(t *testing.T) {
	_, err := domainpkg.NewUser(testDomainUserJohnName, testDomainUserJohnEmail, testDomainUserPassword, domainpkg.UserRole("visitor"))
	if err != domainpkg.ErrUserRoleInvalid {
		t.Fatalf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestUpdateProfileNormalizesValues(t *testing.T) {
	user, err := domainpkg.NewUser(testDomainUserJohnName, testDomainUserJohnEmail, testDomainUserPassword, domainpkg.RoleMonitor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = user.UpdateProfile(" "+testDomainUserJaneName+" ", " Jane@Example.com ", domainpkg.RoleAdmin)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != testDomainUserJaneName {
		t.Fatalf("expected normalized name, got %q", user.Name)
	}
	if user.Email != testDomainUserJaneEmail {
		t.Fatalf("expected normalized email, got %q", user.Email)
	}
	if user.GlobalRole != domainpkg.RoleAdmin {
		t.Fatalf("expected role %q, got %q", domainpkg.RoleAdmin, user.GlobalRole)
	}
}

func TestNewUserRejectsInvalidEmail(t *testing.T) {
	_, err := domainpkg.NewUser(testDomainUserJohnName, "invalid-email", testDomainUserPassword, domainpkg.RoleProfessor)
	if !errors.Is(err, domainpkg.ErrUserEmailInvalid) {
		t.Fatalf("expected ErrUserEmailInvalid, got %v", err)
	}
}

func TestNewUserRejectsShortPassword(t *testing.T) {
	_, err := domainpkg.NewUser(testDomainUserJohnName, testDomainUserJohnEmail, "short", domainpkg.RoleProfessor)
	if !errors.Is(err, domainpkg.ErrUserPasswordTooShort) {
		t.Fatalf("expected ErrUserPasswordTooShort, got %v", err)
	}
}

func TestUpdateProfileRejectsInvalidRole(t *testing.T) {
	user, err := domainpkg.NewUser(testDomainUserJohnName, testDomainUserJohnEmail, testDomainUserPassword, domainpkg.RoleProfessor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = user.UpdateProfile(testDomainUserJohnName, testDomainUserJohnEmail, domainpkg.UserRole("guest"))
	if !errors.Is(err, domainpkg.ErrUserRoleInvalid) {
		t.Fatalf("expected ErrUserRoleInvalid, got %v", err)
	}
}

func TestUpdateProfileRejectsInvalidNameWithoutMutatingUser(t *testing.T) {
	user, err := domainpkg.NewUser(testDomainUserJohnName, testDomainUserJohnEmail, testDomainUserPassword, domainpkg.RoleProfessor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = user.UpdateProfile("  ", testDomainUserJaneEmail, domainpkg.RoleAdmin)
	if !errors.Is(err, domainpkg.ErrUserNameRequired) {
		t.Fatalf("expected ErrUserNameRequired, got %v", err)
	}
	if user.Name != testDomainUserJohnName {
		t.Fatalf("expected original name to remain unchanged, got %q", user.Name)
	}
	if user.Email != testDomainUserJohnEmail {
		t.Fatalf("expected original email to remain unchanged, got %q", user.Email)
	}
}

func TestUpdateProfileRejectsInvalidEmailWithoutMutatingUser(t *testing.T) {
	user, err := domainpkg.NewUser(testDomainUserJohnName, testDomainUserJohnEmail, testDomainUserPassword, domainpkg.RoleProfessor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = user.UpdateProfile(testDomainUserJaneName, "invalid-email", domainpkg.RoleAdmin)
	if !errors.Is(err, domainpkg.ErrUserEmailInvalid) {
		t.Fatalf("expected ErrUserEmailInvalid, got %v", err)
	}
	if user.Name != testDomainUserJohnName {
		t.Fatalf("expected original name to remain unchanged, got %q", user.Name)
	}
	if user.Email != testDomainUserJohnEmail {
		t.Fatalf("expected original email to remain unchanged, got %q", user.Email)
	}
}

func TestNormalizeEmail(t *testing.T) {
	normalized := domainpkg.NormalizeEmail(" John.Doe@Example.com ")
	if normalized != "john.doe@example.com" {
		t.Fatalf("expected normalized email, got %q", normalized)
	}
}

func TestValidateUserNameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 101)
	err := domainpkg.ValidateUserName(longName)
	if !errors.Is(err, domainpkg.ErrUserNameTooLong) {
		t.Fatalf("expected ErrUserNameTooLong, got %v", err)
	}
}
