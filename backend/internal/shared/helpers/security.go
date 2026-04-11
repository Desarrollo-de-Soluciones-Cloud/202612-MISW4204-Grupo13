package helpers

import "golang.org/x/crypto/bcrypt"

const passwordHashCost = bcrypt.DefaultCost

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
