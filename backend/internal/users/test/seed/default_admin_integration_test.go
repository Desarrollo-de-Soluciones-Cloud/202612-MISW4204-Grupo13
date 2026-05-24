package seed_test

import (
	sharedDB "backend/internal/shared/database/testsupport"
	usersDomain "backend/internal/users/domain"
	usersSeed "backend/internal/users/seed"
	"testing"
)

const testDefaultAdminEmail = "admin@example.com"

func TestSeedDefaultAdminCreatesAdminUser(t *testing.T) {
	db := sharedDB.SetupSQLiteDB(t, &usersDomain.User{})
	t.Setenv("DEFAULT_ADMIN_NAME", "Admin User")
	t.Setenv("DEFAULT_ADMIN_EMAIL", testDefaultAdminEmail)
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "supersecret123")

	if err := usersSeed.SeedDefaultAdmin(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var count int64
	if err := db.Model(&usersDomain.User{}).Where("email = ?", testDefaultAdminEmail).Count(&count).Error; err != nil {
		t.Fatalf("expected count query, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one seeded admin, got %d", count)
	}
}

func TestSeedDefaultAdminIsIdempotent(t *testing.T) {
	db := sharedDB.SetupSQLiteDB(t, &usersDomain.User{})
	t.Setenv("DEFAULT_ADMIN_NAME", "Admin User")
	t.Setenv("DEFAULT_ADMIN_EMAIL", testDefaultAdminEmail)
	t.Setenv("DEFAULT_ADMIN_PASSWORD", "supersecret123")

	if err := usersSeed.SeedDefaultAdmin(); err != nil {
		t.Fatalf("expected first seed to succeed, got %v", err)
	}
	if err := usersSeed.SeedDefaultAdmin(); err != nil {
		t.Fatalf("expected second seed to succeed, got %v", err)
	}

	var count int64
	if err := db.Model(&usersDomain.User{}).Where("email = ?", testDefaultAdminEmail).Count(&count).Error; err != nil {
		t.Fatalf("expected count query, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected idempotent seed to keep one admin, got %d", count)
	}
}
