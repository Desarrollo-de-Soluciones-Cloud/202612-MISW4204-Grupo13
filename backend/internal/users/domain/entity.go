package domain

import (
	"strings"
	"time"
)

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Email      string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password   string    `gorm:"size:255;not null" json:"-"`
	GlobalRole UserRole  `gorm:"size:20;not null" json:"global_role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewUser(name, email, password string, role UserRole) (*User, error) {
	normalizedName := strings.TrimSpace(name)
	normalizedEmail := NormalizeEmail(email)

	if err := ValidateUserName(normalizedName); err != nil {
		return nil, err
	}
	if err := ValidateUserEmail(normalizedEmail); err != nil {
		return nil, err
	}
	if err := ValidateUserPassword(password); err != nil {
		return nil, err
	}
	if err := ValidateUserRole(role); err != nil {
		return nil, err
	}

	return &User{
		Name:       normalizedName,
		Email:      normalizedEmail,
		Password:   password,
		GlobalRole: role,
	}, nil
}

func (u *User) UpdateProfile(name, email string, role UserRole) error {
	normalizedName := strings.TrimSpace(name)
	normalizedEmail := NormalizeEmail(email)

	if err := ValidateUserName(normalizedName); err != nil {
		return err
	}
	if err := ValidateUserEmail(normalizedEmail); err != nil {
		return err
	}
	if err := ValidateUserRole(role); err != nil {
		return err
	}

	u.Name = normalizedName
	u.Email = normalizedEmail
	u.GlobalRole = role

	return nil
}
