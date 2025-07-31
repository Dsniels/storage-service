package app

import (
	"context"
	"log"
	"os"

	"github.com/dsniels/storage-service/internal/handler"
	"github.com/dsniels/storage-service/internal/queue"
	store "github.com/dsniels/storage-service/internal/storage"
	pb "github.com/dsniels/storage-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	BlobHandler handler.IBlobHandler
	FileHandler handler.IFileHandler
	Conn        grpc.ClientConn
	Store       store.IStore
	Queue       *queue.Rabbit
}

func InitServices() *App {
	_ = getLogger()
	var grpClient pb.CursosProtoServiceClient = nil
	blobClient := getBlobClient()
	fileClient := getFileClient()
	blobStore := store.NewBlobStore(blobClient)
	fileStore := store.NewFileStore(fileClient)

	conn, err := grpc.Dial(os.Getenv("GRPC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err == nil {
		grpClient = pb.NewCursosProtoServiceClient(conn)
	} else {
		log.Println(err)
	}

	blobHandler := handler.NewBlobHandler(blobStore, blobStore, grpClient)
	fileHandler := handler.NewFileHandler(fileStore)

	rabbit, err := queue.NewRabbit(blobStore)
	if err != nil {
		log.Println(err)
	} else {
		go rabbit.DeleteFileConsumer(context.Background())
	}

	return &App{
		Store:       blobStore,
		BlobHandler: blobHandler,
		FileHandler: fileHandler,
		Queue:       rabbit,
		Conn:        *conn,
	}
}
