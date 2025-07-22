package app

import (
	"log"
	"os"

	"github.com/dsniels/storage-service/internal/controllers"
	"github.com/dsniels/storage-service/internal/queue"
	store "github.com/dsniels/storage-service/internal/storage"
	pb "github.com/dsniels/storage-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	Controller controllers.IController
	Conn       grpc.ClientConn
	Store      store.IStore
	Queue      *queue.Rabbit
}

func InitServices() *App {
	_ = getLogger()
	azClient := getAzureClient()
	blobStore := store.NewBlobStore(azClient)
	conn, err := grpc.Dial(os.Getenv("GRPC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewCursosProtoServiceClient(conn)
	controllers := controllers.NewController(blobStore, blobStore, client)
	rabbit, err := queue.NewRabbit(blobStore)
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		Store:      blobStore,
		Controller: controllers,
		Queue:      rabbit,
		Conn:       *conn,
	}
}
