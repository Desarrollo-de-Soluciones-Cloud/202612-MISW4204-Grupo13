package domain

type UserReader interface {
	GetByEmail(email string) (*AuthenticatedUserCredentials, error)
}

type TokenManager interface {
	GenerateToken(user *AuthenticatedUser) (*AuthToken, error)
	ParseToken(token string) (*AuthenticatedUser, error)
}
