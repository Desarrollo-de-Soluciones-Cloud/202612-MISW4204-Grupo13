package application

import (
	"backend/internal/users/domain"
	"errors"
)

type CreateUserInput struct {
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Password   string          `json:"password"`
	GlobalRole domain.UserRole `json:"global_role"`
}

type CreateUserOutput struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	GlobalRole domain.UserRole `json:"global_role"`
}

type CreateUser struct {
	repository domain.UserRepository
}

func NewCreateUser(repo domain.UserRepository) *CreateUser {
	return &CreateUser{repository: repo}
}

func (uc *CreateUser) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	normalizedEmail := domain.NormalizeEmail(input.Email)

	existing, err := uc.repository.FindByEmail(normalizedEmail)
	if err == nil && existing != nil {
		return nil, domain.ErrUserEmailAlreadyInUse
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	user, err := domain.NewUser(
		input.Name,
		normalizedEmail,
		input.Password,
		input.GlobalRole,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Create(user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: user.GlobalRole,
	}, nil
}
