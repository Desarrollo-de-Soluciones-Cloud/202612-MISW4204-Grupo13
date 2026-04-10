package application

import (
	"backend/internal/users/domain"
)

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserOutput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUser struct {
	repository domain.UserRepository
}

func NewCreateUser(repo domain.UserRepository) *CreateUser {
	return &CreateUser{repository: repo}
}

func (uc *CreateUser) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, domain.ErrInvalidInput
	}

	existing, _ := uc.repository.FindByEmail(input.Email)
	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		Name: input.Name,
		Email: input.Email,
		Password: input.Password,
	}

	if err := uc.repository.Create(user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
	}, nil
}