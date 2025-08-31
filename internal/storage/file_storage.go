package store

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

var dirName = "files"

type FileStore struct {
	*share.Client
}

func (f *FileStore) UploadBlob(ctx context.Context, filename string, content io.Reader, contentType string) (*string, error) {
	return nil, nil
}
func (f *FileStore) UploadFile(ctx context.Context, filename string, content []byte, contentType string) (*string, error) {
	fClient := f.getDir(dirName).NewFileClient(filename)
	size := len(content)
	if _, err := fClient.Create(ctx, int64(size), nil); err != nil {
		return nil, err
	}
	if err := fClient.UploadBuffer(ctx, content, &file.UploadBufferOptions{
		Concurrency: 10,
		ChunkSize:   5 * 1024 * 1024,
	}); err != nil {
		return nil, err
	}
	url := fClient.URL()
	slog.Info("Url: ", url)
	return &url, nil
}

func (f *FileStore) GetFileIdFromURL(ctx context.Context, URL string) (*string, error) {

	return nil, nil
}

func (f *FileStore) GetFiles(ctx context.Context, dirName string, prefix string) (*[]string, error) {

	var files []string
	dir := f.getDir(dirName)
	pages := dir.NewListFilesAndDirectoriesPager(nil)

	for pages.More() {
		res, err := pages.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, file := range res.Segment.Files {
			fmt.Printf("File name: %s", *file.Name)
			files = append(files, *file.Name)
		}
	}

	return &files, nil
}

func (f *FileStore) DeleteFile(ctx context.Context, filename string, dirName string) error {
	fclient := f.getDir(dirName).NewFileClient(filename)
	_, err := fclient.Delete(ctx, nil)
	return err
}

func (f *FileStore) getDir(name string) *directory.Client {
	dir := f.NewDirectoryClient(name)
	dir.Create(context.Background(), nil)
	return dir
}

func NewFileStore(client *share.Client) *FileStore {
	if _, err := client.Create(context.Background(), nil); err != nil {
		slog.Error("Error create share client", err)
	}

	return &FileStore{
		Client: client,
	}
}
