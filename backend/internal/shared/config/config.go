package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                 string
	DBHost               string
	DBPort               string
	DBUnixSocket         string
	CloudSQLConnectionName string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	JWTExpirationMinutes int
	GCPProjectID         string
	GCPLocation          string
	VertexAIModel        string
	GCSBucketName        string
	GCSReportsPrefix     string
	GCSAttachmentsPrefix string
	CORSAllowedOrigins   []string
	ReportsPubSubTopic   string
	PubSubPushAuthToken  string
}

func Load() *Config {
	return &Config{
		Port:                 getEnv("PORT", "8080"),
		DBHost:               getEnv("DB_HOST", ""),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUnixSocket:         getEnv("DB_UNIX_SOCKET", ""),
		CloudSQLConnectionName: getEnv("CLOUD_SQL_CONNECTION_NAME", ""),
		DBUser:               getEnv("DB_USER", ""),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		DBName:               getEnv("DB_NAME", ""),
		JWTSecret:            getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpirationMinutes: getEnvAsInt("JWT_EXPIRATION_MINUTES", 60),
		GCPProjectID:         getEnv("GCP_PROJECT_ID", ""),
		GCPLocation:          getEnv("GCP_LOCATION", "us-central1"),
		VertexAIModel:        getEnv("VERTEX_AI_MODEL", "gemini-2.5-flash-lite"),
		GCSBucketName:        getEnv("GCS_BUCKET_NAME", ""),
		GCSReportsPrefix:     getEnv("GCS_REPORTS_PREFIX", "reports"),
		GCSAttachmentsPrefix: getEnv("GCS_ATTACHMENTS_PREFIX", "attachments"),
		CORSAllowedOrigins:   getEnvAsCSV("CORS_ALLOWED_ORIGINS"),
		ReportsPubSubTopic:   getEnv("REPORTS_PUBSUB_TOPIC", ""),
		PubSubPushAuthToken:  getEnv("PUBSUB_PUSH_AUTH_TOKEN", ""),
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

func getEnvAsCSV(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		result = append(result, trimmed)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
