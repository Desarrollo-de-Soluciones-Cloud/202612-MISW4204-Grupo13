package application

import (
	"backend/internal/users/domain"
	"backend/internal/users/infrastructure"
)

type ListUsersOutput struct {
    Users []UserDTO `json:"users"`
}

type UserDTO struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type ListUsers struct {
    repository domain.UserRepository
}

func NewListUsers(repo *infrastructure.UserRepository) *ListUsers {
    return &ListUsers{repository: repo}
}

func (uc *ListUsers) Execute() (*ListUsersOutput, error) {
    users, err := uc.repository.FindAll()
    if err != nil {
        return nil, err
    }
    result := make([]UserDTO, len(users))
    for i, u := range users {
        result[i] = UserDTO{
            ID:    u.ID,
            Name:  u.Name,
            Email: u.Email,
        }
    }
    return &ListUsersOutput{Users: result}, nil
}