package database

import (
	"fmt"
	"log"
	"time"

	"backend/internal/shared/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := buildPostgresDSN(cfg)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return nil
}

func buildPostgresDSN(cfg *config.Config) string {
	socketPath := resolveDBSocketPath(cfg)
	if socketPath != "" {
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=disable",
			socketPath,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
		)
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
}

func resolveDBSocketPath(cfg *config.Config) string {
	if cfg.DBUnixSocket != "" {
		return cfg.DBUnixSocket
	}

	if cfg.CloudSQLConnectionName != "" {
		return "/cloudsql/" + cfg.CloudSQLConnectionName
	}

	return ""
}
