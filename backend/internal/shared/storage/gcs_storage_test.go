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
