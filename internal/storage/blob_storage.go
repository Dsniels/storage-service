package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	exceptions "github.com/dsniels/storage-service/internal/Exceptions"
)

type BlobStore struct {
	*azblob.Client
}

var defaultName string = "temp"

func (s *BlobStore) GetStream(ctx context.Context, filename string) (io.ReadSeeker, error) {
	container, err := s.getContainer(ctx, defaultName)
	if err != nil {
		return nil, err
	}
	blobClient := container.NewBlobClient(filename)
	props, err := blobClient.GetProperties(ctx, &blob.GetPropertiesOptions{})
	if err != nil {
		return nil, err
	}
	return &BlobReader{Size: *props.ContentLength, Client: blobClient, Pos: 0, Ctx: ctx}, nil
}

func (s *BlobStore) UploadFile(ctx context.Context, filename string, content []byte, contentType string) (*string, error) {
	opts := &azblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
		BlockSize: int64(8 * 1024 * 1024),
	}

	_, err := s.getContainer(ctx, defaultName)
	if err != nil {
		return nil, err
	}
	filename = strings.ReplaceAll(filename, " ", "")

	_, err = s.UploadBuffer(ctx, defaultName, filename, content, opts)
	if err != nil {
		return nil, err
	}

	endpoint := s.URL()

	url := fmt.Sprintf("%s%s/%s", endpoint, defaultName, filename)

	return &url, nil
}

func (s *BlobStore) GetBlobIdFromURL(ctx context.Context, URL string) (*string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	arr := strings.Split(u.Path, "/")
	last := arr[len(arr)-1]

	return &last, nil
}
func (s *BlobStore) GetFiles(ctx context.Context, containerName string, prefix string) (*[]string, error) {
	var blobs []string
	container, err := s.getContainer(ctx, defaultName)
	if err != nil {
		exceptions.ThrowInternalServerError("Couldnt get the container")

	}
	pager := container.NewListBlobsFlatPager(&azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, blob := range response.Segment.BlobItems {
			name := *blob.Name
			blobs = append(blobs, name)
		}
	}
	return &blobs, nil
}

func (s *BlobStore) ensureContainer(ctx context.Context, containerName string) error {

	_, err := s.CreateContainer(ctx, containerName, &azblob.CreateContainerOptions{})

	if err != nil {
		if bloberror.HasCode(err, bloberror.ContainerAlreadyExists) {
			return nil
		}
		return err
	}
	return nil
}

func (s *BlobStore) getContainer(ctx context.Context, containerName string) (*container.Client, error) {

	if containerName == "" {
		containerName = defaultName
	}

	if err := s.ensureContainer(ctx, containerName); err != nil {
		return nil, err
	}

	container := s.ServiceClient().NewContainerClient(containerName)

	return container, nil

}

func (s *BlobStore) DeleteFile(ctx context.Context, filename string, containerName string) error {
	if containerName == "" {
		containerName = defaultName

	}

	container, err := s.getContainer(ctx, containerName)
	if err != nil {
		return err
	}

	blobClient := container.NewBlobClient(filename)

	_, err = blobClient.Delete(ctx, &blob.DeleteOptions{
		DeleteSnapshots: to.Ptr(blob.DeleteSnapshotsOptionTypeInclude),
	})

	if bloberror.HasCode(err, bloberror.BlobNotFound) {
		exceptions.ThrowNotFound()
	}

	return exceptions.ErrorBadRequest

}

func NewBlobStore(client *azblob.Client) *BlobStore {
	return &BlobStore{
		Client: client,
	}
}
