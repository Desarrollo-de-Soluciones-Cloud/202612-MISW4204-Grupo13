package seed

import (
	"os"
	"testing"
)

func TestGetRequiredEnvReturnsValue(t *testing.T) {
	const key = "TEST_DEFAULT_ADMIN_VALUE"
	t.Setenv(key, "admin@example.com")

	value, err := getRequiredEnv(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "admin@example.com" {
		t.Fatalf("expected env value, got %q", value)
	}
}

func TestGetRequiredEnvReturnsErrorWhenMissing(t *testing.T) {
	const key = "TEST_DEFAULT_ADMIN_MISSING"
	_ = os.Unsetenv(key)

	value, err := getRequiredEnv(key)
	if err == nil {
		t.Fatalf("expected missing env error")
	}
	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}
	if err.Error() != key+" is required" {
		t.Fatalf("unexpected error %v", err)
	}
}
