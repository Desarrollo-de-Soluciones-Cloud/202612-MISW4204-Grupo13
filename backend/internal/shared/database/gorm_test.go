package database

import (
	sharedConfig "backend/internal/shared/config"
	"testing"
)

const (
	testSocketPath  = "/tmp/socket"
	errExpectedFmt  = "expected %q, got %q"
)

func TestBuildPostgresDSNUsesSocketWhenPresent(t *testing.T) {
	cfg := &sharedConfig.Config{
		DBUnixSocket: testSocketPath,
		DBUser:       "user",
		DBPassword:   "pass",
		DBName:       "dbname",
	}

	dsn := buildPostgresDSN(cfg)
	expected := "host=/tmp/socket user=user password=pass dbname=dbname sslmode=disable"
	if dsn != expected {
		t.Fatalf(errExpectedFmt, expected, dsn)
	}
}

func TestBuildPostgresDSNUsesCloudSQLSocketFallback(t *testing.T) {
	cfg := &sharedConfig.Config{
		CloudSQLConnectionName: "project:region:instance",
		DBUser:                 "user",
		DBPassword:             "pass",
		DBName:                 "dbname",
	}

	dsn := buildPostgresDSN(cfg)
	expected := "host=/cloudsql/project:region:instance user=user password=pass dbname=dbname sslmode=disable"
	if dsn != expected {
		t.Fatalf(errExpectedFmt, expected, dsn)
	}
}

func TestBuildPostgresDSNUsesHostAndPort(t *testing.T) {
	cfg := &sharedConfig.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "user",
		DBPassword: "pass",
		DBName:     "dbname",
	}

	dsn := buildPostgresDSN(cfg)
	expected := "host=localhost user=user password=pass dbname=dbname port=5432 sslmode=disable"
	if dsn != expected {
		t.Fatalf(errExpectedFmt, expected, dsn)
	}
}

func TestResolveDBSocketPath(t *testing.T) {
	cfg := &sharedConfig.Config{DBUnixSocket: testSocketPath, CloudSQLConnectionName: "ignored"}
	if path := resolveDBSocketPath(cfg); path != testSocketPath {
		t.Fatalf("expected explicit socket path, got %q", path)
	}

	cfg = &sharedConfig.Config{CloudSQLConnectionName: "project:region:instance"}
	if path := resolveDBSocketPath(cfg); path != "/cloudsql/project:region:instance" {
		t.Fatalf("expected cloudsql socket path, got %q", path)
	}

	cfg = &sharedConfig.Config{}
	if path := resolveDBSocketPath(cfg); path != "" {
		t.Fatalf("expected empty socket path, got %q", path)
	}
}
