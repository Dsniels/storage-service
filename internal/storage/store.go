package store

import (
	"context"
	"io"
)

type IStore interface {
	GetFileIdFromURL(ctx context.Context, URL string) (*string, error)
	UploadBlob(ctx context.Context, filename string, content io.Reader, contentType string) (*string, error)
	UploadFile(ctx context.Context, filename string, content []byte, contentType string) (*string, error)
	GetFiles(context.Context, string, string) (*[]string, error)
	DeleteFile(context.Context, string, string) error
}

type IStream interface {
	GetStream(ctx context.Context, filename string) (io.ReadSeeker, error)
}
