package infrastructure

import (
	"backend/internal/auth/domain"
	usersApplication "backend/internal/users/application"
)

type UserReader struct {
	getUserByEmail *usersApplication.GetUserByEmail
}

func NewUserReader(getUserByEmail *usersApplication.GetUserByEmail) *UserReader {
	return &UserReader{getUserByEmail: getUserByEmail}
}

func (r *UserReader) GetByEmail(email string) (*domain.AuthenticatedUserCredentials, error) {
	output, err := r.getUserByEmail.Execute(usersApplication.GetUserByEmailInput{Email: email})
	if err != nil {
		return nil, err
	}

	return &domain.AuthenticatedUserCredentials{
		AuthenticatedUser: domain.AuthenticatedUser{
			ID:         output.ID,
			Name:       output.Name,
			Email:      output.Email,
			GlobalRole: output.GlobalRole,
		},
		Password: output.Password,
	}, nil
}
