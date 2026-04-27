package storage

import (
	"context"
	"errors"
	"io"

	"cloud.google.com/go/storage"
)

type GCSStorage struct {
	client     *storage.Client
	bucketName string
}

func NewGCSStorage(ctx context.Context, bucketName string) (*GCSStorage, error) {
	if bucketName == "" {
		return nil, errors.New("gcs bucket name is required")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCSStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *GCSStorage) Upload(ctx context.Context, objectName string, reader io.Reader, contentType string) error {
	if objectName == "" {
		return errors.New("gcs object name is required")
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w := s.client.Bucket(s.bucketName).Object(objectName).NewWriter(ctx)
	w.ContentType = contentType

	if _, err := io.Copy(w, reader); err != nil {
		_ = w.Close()
		return err
	}

	return w.Close()
}

func (s *GCSStorage) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	if objectName == "" {
		return nil, errors.New("gcs object name is required")
	}

	return s.client.Bucket(s.bucketName).Object(objectName).NewReader(ctx)
}

func (s *GCSStorage) Close() error {
	return s.client.Close()
}