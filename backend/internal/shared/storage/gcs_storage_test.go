package storage

import (
	"context"
	"strings"
	"testing"
)

func TestNewGCSStorageRejectsEmptyBucketName(t *testing.T) {
	_, err := NewGCSStorage(context.Background(), "")
	if err == nil || err.Error() != "gcs bucket name is required" {
		t.Fatalf("expected bucket name required error, got %v", err)
	}
}

func TestGCSStorageRejectsEmptyObjectName(t *testing.T) {
	storage := &GCSStorage{}

	if err := storage.Upload(context.Background(), "", strings.NewReader("data"), "text/plain"); err == nil || err.Error() != errGCSObjectNameRequired {
		t.Fatalf("expected object name required on upload, got %v", err)
	}

	if _, err := storage.Download(context.Background(), ""); err == nil || err.Error() != errGCSObjectNameRequired {
		t.Fatalf("expected object name required on download, got %v", err)
	}

	if err := storage.Delete(context.Background(), ""); err == nil || err.Error() != errGCSObjectNameRequired {
		t.Fatalf("expected object name required on delete, got %v", err)
	}
}

func TestNewGCSStorageWithEmulatorAndOperationsReachClientPath(t *testing.T) {
	t.Setenv("STORAGE_EMULATOR_HOST", "http://127.0.0.1:9090")

	storage, err := NewGCSStorage(context.Background(), "bucket")
	if err != nil {
		t.Fatalf("expected storage creation with emulator host, got %v", err)
	}
	defer func() { _ = storage.Close() }()

	if err := storage.Upload(context.Background(), "reports/file.txt", strings.NewReader("hello"), ""); err == nil {
		t.Fatalf("expected upload to fail without emulator server")
	}
	if _, err := storage.Download(context.Background(), "reports/file.txt"); err == nil {
		t.Fatalf("expected download to fail without emulator server")
	}
	if err := storage.Delete(context.Background(), "reports/file.txt"); err == nil {
		t.Fatalf("expected delete to fail without emulator server")
	}
}
