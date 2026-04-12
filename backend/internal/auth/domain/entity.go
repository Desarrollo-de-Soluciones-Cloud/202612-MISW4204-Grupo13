package domain

import usersDomain "backend/internal/users/domain"

type AuthenticatedUser struct {
	ID         uint                 `json:"id"`
	Name       string               `json:"name"`
	Email      string               `json:"email"`
	GlobalRole usersDomain.UserRole `json:"global_role"`
}

type AuthenticatedUserCredentials struct {
	AuthenticatedUser
	Password string `json:"-"`
}

type AuthToken struct {
	AccessToken string    `json:"access_token"`
	TokenType   TokenType `json:"token_type"`
	ExpiresIn   int64     `json:"expires_in"`
}
