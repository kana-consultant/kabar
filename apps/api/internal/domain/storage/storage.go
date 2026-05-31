package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) (string, error)
	GetURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	Delete(ctx context.Context, objectName string) error
	List(ctx context.Context, prefix string) ([]string, error)
}
