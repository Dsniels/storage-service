package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
)

type Store struct {
	client *azblob.Client
	logger *log.Logger
}

var container = "temp"

type IStore interface {
	GetBlobIdFromURL(ctx context.Context, URL string) (*string, error)
	UploadFile(ctx context.Context, filename string, content []byte, contentType string) (*string, error)
	GetFileStream(ctx context.Context, filename string) (io.ReadCloser, error)
	GetFilesfromContainer(ctx context.Context, container string)
	EnsureContainer(ctx context.Context, containerName string) error
}

func NewStore(client *azblob.Client, logger *log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) GetFileStream(ctx context.Context, filename string, offSet int64) (io.ReadSeeker, error) {
	blobClient := s.client.ServiceClient().NewContainerClient(container).NewBlobClient(filename)
	props, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &BlobReader{Size: *props.ContentLength, Client: blobClient, Pos: 0, Ctx: ctx}, nil
}

func (s *Store) UploadFile(ctx context.Context, filename string, content []byte, contentType string) (*string, error) {
	container := "temp"
	opts := &azblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
		BlockSize: int64(8 * 1024 * 1024),
	}

	err := s.EnsureContainer(ctx, container)
	if err != nil {
		return nil, err
	}
	filename = strings.ReplaceAll(filename, " ", "")

	_, err = s.client.UploadBuffer(ctx, container, filename, content, opts)
	if err != nil {
		return nil, err
	}

	endpoint := s.client.URL()
	s.logger.Println(endpoint)

	url := fmt.Sprintf("%s%s/%s", endpoint, container, filename)
	s.logger.Println(url)

	return &url, nil
}

func (s *Store) GetBlobIdFromURL(ctx context.Context, URL string) (*string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	arr := strings.Split(u.Path, "/")
	last := arr[len(arr)-1]

	return &last, nil
}

func (s *Store) EnsureContainer(ctx context.Context, containerName string) error {

	_, err := s.client.CreateContainer(ctx, containerName, &azblob.CreateContainerOptions{})

	if err != nil {
		if bloberror.HasCode(err, bloberror.ContainerAlreadyExists) {
			return nil
		}
		return err
	}
	return nil
}
