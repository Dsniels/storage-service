package app

import (
	"github.com/dsniels/storage-service/internal/controllers"
	"github.com/dsniels/storage-service/internal/storage"
)

type App struct {
	Controller controllers.IController
	Store      storage.IStore
}

func InitServices() *App {
	_ = getLogger()
	azClient := getAzureClient()
	blobStore := storage.NewBlobStore(azClient)
	controllers := controllers.NewController(blobStore, blobStore)

	return &App{
		Store:      blobStore,
		Controller: controllers,
	}
}
