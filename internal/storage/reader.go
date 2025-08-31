package store

import (
	"context"
	"io"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

type BlobReader struct {
	Ctx         context.Context
	Pos         int64
	Size        int64
	buffer      []byte
	bufferStart int64
	*blob.Client
}

func (b *BlobReader) Read(p []byte) (int, error) {
	if b.Pos >= b.Size {
		return 0, io.EOF
	}

	if b.Pos >= b.bufferStart && b.Pos < b.bufferStart+int64(len(b.buffer)) {
		start := b.Pos - b.bufferStart
		n := copy(p, b.buffer[start:])
		b.Pos += int64(n)
		return n, nil
	}

	chunkSize := int64(len(p))
	if remaining := b.Size - b.Pos; remaining < chunkSize {
		chunkSize = remaining
	}

	res, err := b.Client.DownloadBuffer(b.Ctx, p, &blob.DownloadBufferOptions{
		Progress: func(bytesTransferred int64) {
			slog.Info("Downloading", slog.Int64("bytes", bytesTransferred))
		},
		Concurrency: 6,
		BlockSize:   2 * 1024 * 1024,
		Range:       azblob.HTTPRange{Count: chunkSize, Offset: b.Pos},
	})
	if err != nil {
		slog.Error("Download Stream: ", err)
		return 0, err
	}
	copy(b.buffer, p)
	b.bufferStart = b.Pos
	b.Pos += res
	return int(res), nil
}

func (b *BlobReader) Seek(offset int64, whence int) (int64, error) {

	switch whence {
	case io.SeekStart:
		b.Pos = offset
	case io.SeekEnd:
		b.Pos = b.Size + offset
	case io.SeekCurrent:
		b.Pos += offset
	}

	if b.Pos < 0 {
		b.Pos = 0
	}
	return b.Pos, nil

}
