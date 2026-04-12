package seed

import (
	"backend/internal/users/application"
	"backend/internal/users/domain"
	"backend/internal/users/infrastructure"
	"errors"
	"os"
)

const (
	defaultAdminNameEnv     = "DEFAULT_ADMIN_NAME"
	defaultAdminEmailEnv    = "DEFAULT_ADMIN_EMAIL"
	defaultAdminPasswordEnv = "DEFAULT_ADMIN_PASSWORD"
)

func SeedDefaultAdmin() error {
	repo := infrastructure.NewUserRepository()

	if err := repo.AutoMigrate(); err != nil {
		return err
	}

	defaultAdminName, err := getRequiredEnv(defaultAdminNameEnv)
	if err != nil {
		return err
	}

	defaultAdminEmail, err := getRequiredEnv(defaultAdminEmailEnv)
	if err != nil {
		return err
	}

	defaultAdminPassword, err := getRequiredEnv(defaultAdminPasswordEnv)
	if err != nil {
		return err
	}

	_, err := repo.FindByEmail(domain.NormalizeEmail(defaultAdminEmail))
	if err == nil {
		return nil
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}

	createUser := application.NewCreateUser(repo)
	_, err = createUser.Execute(application.CreateUserInput{
		Name:       defaultAdminName,
		Email:      defaultAdminEmail,
		Password:   defaultAdminPassword,
		GlobalRole: domain.RoleAdmin,
	})

	return err
}

func getRequiredEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	return "", errors.New(key + " is required")
}
