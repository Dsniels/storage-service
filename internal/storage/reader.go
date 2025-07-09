package storage

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

type BlobReader struct {
	Ctx       context.Context
	Pos       int64
	Size      int64
	Reader    io.ReadCloser
	container string
	fileName  string
	*blob.Client
}

func (b *BlobReader) Read(p []byte) (int, error) {
	return n, nil
}

func (b *BlobReader) Seek(offset int64, whence int) (int64, error) {

	var newOffset int64

	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekEnd:
		newOffset = b.Offset + offset
	}

	b.Reader.Close()
	b.Reader = nil

	b.Offset = newOffset
	b.Pos = 0

	return b.Offset, nil

}
