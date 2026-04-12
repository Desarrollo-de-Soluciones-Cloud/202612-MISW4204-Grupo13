package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                 string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	JWTExpirationMinutes int
}

func Load() *Config {
	return &Config{
		Port:                 getEnv("PORT", "8080"),
		DBHost:               getEnv("DB_HOST", ""),
		DBPort:               getEnv("DB_PORT", ""),
		DBUser:               getEnv("DB_USER", ""),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		DBName:               getEnv("DB_NAME", ""),
		JWTSecret:            getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpirationMinutes: getEnvAsInt("JWT_EXPIRATION_MINUTES", 60),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
