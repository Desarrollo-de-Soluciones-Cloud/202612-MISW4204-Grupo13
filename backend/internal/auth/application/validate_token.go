package application

import "backend/internal/auth/domain"

type ValidateTokenInput struct {
	Token string `json:"token"`
}

type ValidateTokenOutput struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	GlobalRole string          `json:"global_role"`
	TokenType  domain.TokenType `json:"token_type"`
}

type ValidateToken struct {
	tokenManager domain.TokenManager
}

func NewValidateToken(tokenManager domain.TokenManager) *ValidateToken {
	return &ValidateToken{tokenManager: tokenManager}
}

func (uc *ValidateToken) Execute(input ValidateTokenInput) (*ValidateTokenOutput, error) {
	if err := domain.ValidateTokenString(input.Token); err != nil {
		return nil, err
	}

	user, err := uc.tokenManager.ParseToken(input.Token)
	if err != nil {
		return nil, err
	}

	return &ValidateTokenOutput{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: string(user.GlobalRole),
		TokenType:  domain.TokenTypeBearer,
	}, nil
}
