package infrastructure

import (
	"backend/internal/auth/domain"
	usersDomain "backend/internal/users/domain"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type TokenManager struct {
	secret            []byte
	expirationMinutes int
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type tokenClaims struct {
	Subject   uint   `json:"sub"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	GlobalRole string `json:"global_role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func NewTokenManager(secret string, expirationMinutes int) *TokenManager {
	if expirationMinutes <= 0 {
		expirationMinutes = 60
	}

	return &TokenManager{
		secret:            []byte(secret),
		expirationMinutes: expirationMinutes,
	}
}

func (m *TokenManager) GenerateToken(user *domain.AuthenticatedUser) (*domain.AuthToken, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(m.expirationMinutes) * time.Minute)

	headerBytes, _ := json.Marshal(tokenHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	})

	claimsBytes, _ := json.Marshal(tokenClaims{
		Subject:    user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: string(user.GlobalRole),
		IssuedAt:   now.Unix(),
		ExpiresAt:  expiresAt.Unix(),
	})

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)
	signingInput := encodedHeader + "." + encodedClaims
	signature := m.sign(signingInput)
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return &domain.AuthToken{
		AccessToken: signingInput + "." + encodedSignature,
		TokenType:   domain.TokenTypeBearer,
		ExpiresIn:   int64(time.Duration(m.expirationMinutes) * time.Minute / time.Second),
	}, nil
}

func (m *TokenManager) ParseToken(token string) (*domain.AuthenticatedUser, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, domain.ErrAuthTokenInvalid
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := m.sign(signingInput)
	receivedSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, domain.ErrAuthTokenInvalid
	}

	if subtle.ConstantTimeCompare(receivedSignature, expectedSignature) != 1 {
		return nil, domain.ErrAuthTokenInvalid
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, domain.ErrAuthTokenInvalid
	}

	var claims tokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, domain.ErrAuthTokenInvalid
	}

	if time.Now().UTC().Unix() > claims.ExpiresAt {
		return nil, domain.ErrAuthTokenExpired
	}

	role := usersDomain.UserRole(claims.GlobalRole)
	if err := usersDomain.ValidateUserRole(role); err != nil {
		return nil, domain.ErrAuthTokenInvalid
	}

	return &domain.AuthenticatedUser{
		ID:         claims.Subject,
		Name:       claims.Name,
		Email:      claims.Email,
		GlobalRole: role,
	}, nil
}

func (m *TokenManager) sign(payload string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}
