package application

import "backend/internal/users/domain"

type ChangeUserRoleInput struct {
	ID         uint
	GlobalRole domain.UserRole
}

type ChangeUserRoleOutput struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	GlobalRole domain.UserRole `json:"global_role"`
}

type ChangeUserRole struct {
	repository domain.UserRepository
}

func NewChangeUserRole(repo domain.UserRepository) *ChangeUserRole {
	return &ChangeUserRole{repository: repo}
}

func (uc *ChangeUserRole) Execute(input ChangeUserRoleInput) (*ChangeUserRoleOutput, error) {
	if err := domain.ValidateUserRole(input.GlobalRole); err != nil {
		return nil, err
	}

	user, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	user.GlobalRole = input.GlobalRole
	if err := uc.repository.Update(user); err != nil {
		return nil, err
	}

	return &ChangeUserRoleOutput{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GlobalRole: user.GlobalRole,
	}, nil
}
