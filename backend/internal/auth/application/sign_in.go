package application

import (
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/auth/domain"
	usersDomain "backend/internal/users/domain"
	"errors"
)

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInOutput struct {
	AccessToken string               `json:"access_token"`
	TokenType   domain.TokenType     `json:"token_type"`
	ExpiresIn   int64                `json:"expires_in"`
	User        AuthenticatedUserDTO `json:"user"`
}

type AuthenticatedUserDTO struct {
	ID         uint                 `json:"id"`
	Name       string               `json:"name"`
	Email      string               `json:"email"`
	GlobalRole usersDomain.UserRole `json:"global_role"`
}

type SignIn struct {
	userReader   domain.UserReader
	tokenManager domain.TokenManager
}

func NewSignIn(userReader domain.UserReader, tokenManager domain.TokenManager) *SignIn {
	return &SignIn{
		userReader:   userReader,
		tokenManager: tokenManager,
	}
}

func (uc *SignIn) Execute(input SignInInput) (*SignInOutput, error) {
	normalizedEmail := domain.NormalizeAuthEmail(input.Email)

	if err := domain.ValidateAuthEmail(normalizedEmail); err != nil {
		return nil, err
	}
	if err := domain.ValidateAuthPassword(input.Password); err != nil {
		return nil, err
	}

	user, err := uc.userReader.GetByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, usersDomain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := sharedHelpers.ComparePassword(user.Password, input.Password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := uc.tokenManager.GenerateToken(&domain.AuthenticatedUser{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: user.GlobalRole,
	})
	if err != nil {
		return nil, err
	}

	return &SignInOutput{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   token.ExpiresIn,
		User: AuthenticatedUserDTO{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			GlobalRole: user.GlobalRole,
		},
	}, nil
}
