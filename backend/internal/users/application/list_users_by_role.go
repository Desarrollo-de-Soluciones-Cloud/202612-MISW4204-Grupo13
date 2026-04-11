package application

import "backend/internal/users/domain"

type ListUsersByRoleInput struct {
	GlobalRole domain.UserRole
}

type ListUsersByRoleOutput struct {
	Users []UserDTO `json:"users"`
}

type ListUsersByRole struct {
	repository domain.UserRepository
}

func NewListUsersByRole(repo domain.UserRepository) *ListUsersByRole {
	return &ListUsersByRole{repository: repo}
}

func (uc *ListUsersByRole) Execute(input ListUsersByRoleInput) (*ListUsersByRoleOutput, error) {
	if err := domain.ValidateUserRole(input.GlobalRole); err != nil {
		return nil, err
	}

	users, err := uc.repository.FindAllByRole(input.GlobalRole)
	if err != nil {
		return nil, err
	}

	result := make([]UserDTO, len(users))
	for i, u := range users {
		result[i] = UserDTO{
			ID:         u.ID,
			Name:       u.Name,
			Email:      u.Email,
			GlobalRole: u.GlobalRole,
		}
	}

	return &ListUsersByRoleOutput{Users: result}, nil
}
