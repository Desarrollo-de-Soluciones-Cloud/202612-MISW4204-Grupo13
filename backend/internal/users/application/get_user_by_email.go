package application

import "backend/internal/users/domain"

type GetUserByEmailInput struct {
	Email string
}

type GetUserByEmailOutput struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	GlobalRole domain.UserRole `json:"global_role"`
	Password   string          `json:"-"`
}

type GetUserByEmail struct {
	repository domain.UserRepository
}

func NewGetUserByEmail(repo domain.UserRepository) *GetUserByEmail {
	return &GetUserByEmail{repository: repo}
}

func (uc *GetUserByEmail) Execute(input GetUserByEmailInput) (*GetUserByEmailOutput, error) {
	if err := domain.ValidateUserEmail(input.Email); err != nil {
		return nil, err
	}

	user, err := uc.repository.FindByEmail(domain.NormalizeEmail(input.Email))
	if err != nil {
		return nil, err
	}

	return &GetUserByEmailOutput{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: user.GlobalRole,
		Password:   user.Password,
	}, nil
}
