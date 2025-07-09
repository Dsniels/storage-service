package app

import (
	"github.com/dsniels/storage-service/internal/controllers"
	"github.com/dsniels/storage-service/internal/storage"
)

type App struct {
	Controller controllers.IController
	Store      storage.IStore
}

func NewApp(store storage.IStore, controller controllers.IController) *App {
	return &App{
		Store:      store,
		Controller: controller,
	}
}
