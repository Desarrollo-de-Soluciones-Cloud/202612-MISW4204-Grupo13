package application

import "backend/internal/users/domain"

type GetUserByIDInput struct {
	ID uint
}

type GetUserByIDOutput struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	GlobalRole domain.UserRole `json:"global_role"`
}

type GetUserByID struct {
	repository domain.UserRepository
}

func NewGetUserByID(repo domain.UserRepository) *GetUserByID {
	return &GetUserByID{repository: repo}
}

func (uc *GetUserByID) Execute(input GetUserByIDInput) (*GetUserByIDOutput, error) {
	user, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserByIDOutput{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: user.GlobalRole,
	}, nil
}
