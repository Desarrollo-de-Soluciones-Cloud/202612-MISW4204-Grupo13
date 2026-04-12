package infrastructure

import (
	authDomain "backend/internal/auth/domain"
	authInfrastructure "backend/internal/auth/infrastructure"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	usersDomain "backend/internal/users/domain"
)

func TestNewTokenManagerUsesDefaultExpiration(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 0)
	if manager == nil {
		t.Fatal("expected token manager")
	}
}

func TestGenerateAndParseTokenSuccess(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)
	user := &authDomain.AuthenticatedUser{
		ID:         1,
		Name:       "John Doe",
		Email:      "john@example.com",
		GlobalRole: usersDomain.RoleProfessor,
	}

	token, err := manager.GenerateToken(user)
	if err != nil {
		t.Fatalf("expected token, got %v", err)
	}

	parsedUser, err := manager.ParseToken(token.AccessToken)
	if err != nil {
		t.Fatalf("expected parsed user, got %v", err)
	}
	if parsedUser.Email != user.Email {
		t.Fatalf("expected email %q, got %q", user.Email, parsedUser.Email)
	}
}

func TestParseTokenRejectsInvalidStructure(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)

	_, err := manager.ParseToken("invalid-token")
	if !errors.Is(err, authDomain.ErrAuthTokenInvalid) {
		t.Fatalf("expected ErrAuthTokenInvalid, got %v", err)
	}
}

func TestParseTokenRejectsInvalidSignature(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)
	user := &authDomain.AuthenticatedUser{
		ID:         1,
		Name:       "John Doe",
		Email:      "john@example.com",
		GlobalRole: usersDomain.RoleProfessor,
	}

	token, err := manager.GenerateToken(user)
	if err != nil {
		t.Fatalf("expected token, got %v", err)
	}

	invalidToken := token.AccessToken[:len(token.AccessToken)-1] + "x"
	_, err = manager.ParseToken(invalidToken)
	if !errors.Is(err, authDomain.ErrAuthTokenInvalid) {
		t.Fatalf("expected ErrAuthTokenInvalid, got %v", err)
	}
}

func TestParseTokenRejectsExpiredToken(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)
	expiredToken := buildSignedToken(t, "secret", tokenClaims{
		Subject:    1,
		Name:       "Expired User",
		Email:      "expired@example.com",
		GlobalRole: string(usersDomain.RoleProfessor),
		IssuedAt:   time.Now().UTC().Add(-2 * time.Hour).Unix(),
		ExpiresAt:  time.Now().UTC().Add(-time.Hour).Unix(),
	})

	_, err := manager.ParseToken(expiredToken)
	if !errors.Is(err, authDomain.ErrAuthTokenExpired) {
		t.Fatalf("expected ErrAuthTokenExpired, got %v", err)
	}
}

func TestParseTokenRejectsInvalidRole(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)
	invalidRoleToken := buildSignedToken(t, "secret", tokenClaims{
		Subject:    1,
		Name:       "Role User",
		Email:      "role@example.com",
		GlobalRole: "guest",
		IssuedAt:   time.Now().UTC().Unix(),
		ExpiresAt:  time.Now().UTC().Add(time.Hour).Unix(),
	})

	_, err := manager.ParseToken(invalidRoleToken)
	if !errors.Is(err, authDomain.ErrAuthTokenInvalid) {
		t.Fatalf("expected ErrAuthTokenInvalid, got %v", err)
	}
}

func TestParseTokenRejectsInvalidBase64Signature(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)

	_, err := manager.ParseToken("a.b.invalid@@")
	if !errors.Is(err, authDomain.ErrAuthTokenInvalid) {
		t.Fatalf("expected ErrAuthTokenInvalid, got %v", err)
	}
}

func TestParseTokenRejectsInvalidClaimsPayload(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)
	token := buildMalformedClaimsToken(t, "secret")

	_, err := manager.ParseToken(token)
	if !errors.Is(err, authDomain.ErrAuthTokenInvalid) {
		t.Fatalf("expected ErrAuthTokenInvalid, got %v", err)
	}
}

func TestParseTokenRejectsInvalidClaimsBase64(t *testing.T) {
	manager := authInfrastructure.NewTokenManager("secret", 10)
	token := buildTokenWithInvalidClaimsBase64(t, "secret")

	_, err := manager.ParseToken(token)
	if !errors.Is(err, authDomain.ErrAuthTokenInvalid) {
		t.Fatalf("expected ErrAuthTokenInvalid, got %v", err)
	}
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type tokenClaims struct {
	Subject    uint   `json:"sub"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GlobalRole string `json:"global_role"`
	IssuedAt   int64  `json:"iat"`
	ExpiresAt  int64  `json:"exp"`
}

func buildSignedToken(t *testing.T, secret string, claims tokenClaims) string {
	t.Helper()

	headerBytes, err := json.Marshal(tokenHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	})
	if err != nil {
		t.Fatalf("expected header marshal, got %v", err)
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("expected claims marshal, got %v", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)
	signingInput := encodedHeader + "." + encodedClaims

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature
}

func buildMalformedClaimsToken(t *testing.T, secret string) string {
	t.Helper()

	headerBytes, err := json.Marshal(tokenHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	})
	if err != nil {
		t.Fatalf("expected header marshal, got %v", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedClaims := base64.RawURLEncoding.EncodeToString([]byte("{"))
	signingInput := encodedHeader + "." + encodedClaims

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature
}

func buildTokenWithInvalidClaimsBase64(t *testing.T, secret string) string {
	t.Helper()

	headerBytes, err := json.Marshal(tokenHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	})
	if err != nil {
		t.Fatalf("expected header marshal, got %v", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	signingInput := encodedHeader + ".invalid@@"

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature
}
