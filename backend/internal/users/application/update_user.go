package application

import (
	"backend/internal/users/domain"
	"errors"
)

type UpdateUserInput struct {
	ID         uint
	Name       string
	Email      string
	GlobalRole domain.UserRole
}

type UpdateUserOutput struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	GlobalRole domain.UserRole `json:"global_role"`
}

type UpdateUser struct {
	repository domain.UserRepository
}

func NewUpdateUser(repo domain.UserRepository) *UpdateUser {
	return &UpdateUser{repository: repo}
}

func (uc *UpdateUser) Execute(input UpdateUserInput) (*UpdateUserOutput, error) {
	user, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	normalizedEmail := domain.NormalizeEmail(input.Email)
	if normalizedEmail != user.Email {
		existing, err := uc.repository.FindByEmail(normalizedEmail)
		if err == nil && existing != nil && existing.ID != user.ID {
			return nil, domain.ErrUserEmailAlreadyInUse
		}
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
	}

	if err := user.UpdateProfile(input.Name, normalizedEmail, input.GlobalRole); err != nil {
		return nil, err
	}

	if err := uc.repository.Update(user); err != nil {
		return nil, err
	}

	return &UpdateUserOutput{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: user.GlobalRole,
	}, nil
}
