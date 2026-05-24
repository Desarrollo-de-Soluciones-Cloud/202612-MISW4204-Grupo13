package config

import (
	"os"
	"testing"
)

func TestLoadUsesDefaultsAndConfiguredValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("JWT_EXPIRATION_MINUTES", "120")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000, https://frontend.example.com")
	t.Setenv("REPORTS_PUBSUB_TOPIC", "weekly-reports")
	t.Setenv("PUBSUB_PUSH_AUTH_TOKEN", "push-token")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Port)
	}
	if cfg.DBHost != "db.internal" {
		t.Fatalf("expected DB host db.internal, got %q", cfg.DBHost)
	}
	if cfg.JWTExpirationMinutes != 120 {
		t.Fatalf("expected JWT expiration 120, got %d", cfg.JWTExpirationMinutes)
	}
	if len(cfg.CORSAllowedOrigins) != 2 {
		t.Fatalf("expected 2 allowed origins, got %d", len(cfg.CORSAllowedOrigins))
	}
	if cfg.ReportsPubSubTopic != "weekly-reports" {
		t.Fatalf("expected pubsub topic weekly-reports, got %q", cfg.ReportsPubSubTopic)
	}
	if cfg.PubSubPushAuthToken != "push-token" {
		t.Fatalf("expected push token push-token, got %q", cfg.PubSubPushAuthToken)
	}
}

func TestLoadFallsBackToDefaultsWhenEnvValuesAreMissingOrInvalid(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("JWT_EXPIRATION_MINUTES", "invalid")
	t.Setenv("CORS_ALLOWED_ORIGINS", " , ")
	t.Setenv("VERTEX_AI_MODEL", "")
	t.Setenv("GCP_LOCATION", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if cfg.DBPort != "5432" {
		t.Fatalf("expected default DB port 5432, got %q", cfg.DBPort)
	}
	if cfg.JWTExpirationMinutes != 60 {
		t.Fatalf("expected default JWT expiration 60, got %d", cfg.JWTExpirationMinutes)
	}
	if cfg.CORSAllowedOrigins != nil {
		t.Fatalf("expected nil CORS allowed origins, got %#v", cfg.CORSAllowedOrigins)
	}
	if cfg.VertexAIModel != "gemini-2.5-flash-lite" {
		t.Fatalf("expected default model, got %q", cfg.VertexAIModel)
	}
	if cfg.GCPLocation != "us-central1" {
		t.Fatalf("expected default GCP location us-central1, got %q", cfg.GCPLocation)
	}
}

func TestGetEnvReturnsDefaultValue(t *testing.T) {
	key := "CONFIG_TEST_UNSET"
	_ = os.Unsetenv(key)

	if value := getEnv(key, "fallback"); value != "fallback" {
		t.Fatalf("expected fallback value, got %q", value)
	}
}
