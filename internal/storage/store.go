package storage

import (
	"context"
	"io"
)

type IStore interface {
	GetBlobIdFromURL(ctx context.Context, URL string) (*string, error)
	UploadFile(ctx context.Context, filename string, content []byte, contentType string) (*string, error)
	GetFiles(ctx context.Context, containerName string, prefix string) (*[]string, error)
	DeleteFile(ctx context.Context, filename string, containerName string) error
}

type IStream interface {
	GetStream(ctx context.Context, filename string) (io.ReadSeeker, error)
}
