//go:build wireinject
// +build wireinject

package provider

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/dsniels/storage-service/internal/app"
	"github.com/dsniels/storage-service/internal/controllers"
	"github.com/dsniels/storage-service/internal/storage"
	"github.com/google/wire"
)

func getAzureClient() *azblob.Client {

	connection := os.Getenv("connection")
	log.Println(connection)
	client, err := azblob.NewClientFromConnectionString(connection, &azblob.ClientOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func getLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

var storeSet = wire.NewSet(storage.NewStore, getAzureClient, getLogger, wire.Bind(new(storage.IStore), new(*storage.Store)))
var controllerSet = wire.NewSet(controllers.NewController, wire.Bind(new(controllers.IController), new(*controllers.Controller)))

func Inject() *app.App {
	wire.Build(app.NewApp, storeSet, controllerSet)
	return nil
}
