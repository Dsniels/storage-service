package router

import (
	"github.com/dsniels/storage-service/internal/app"
	"github.com/dsniels/storage-service/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(app *app.App) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Exception)
	router.Route("/api/store", func(r chi.Router) {
		r.Post("/UploadFile", app.BlobHandler.HandleUploadFile)
		r.Delete("/Delete/{id}", app.BlobHandler.HandleDeleteFile)
		r.Get("/Stream/{id}", app.BlobHandler.HandleStreamFile)
		r.Get("/ListFiles", app.BlobHandler.HandleListFiles)
	})
	router.Route("/api/fileStore", func(r chi.Router) {
		r.Post("/UploadFile", app.FileHandler.HandleUploadFile)
		r.Delete("/Delete/{id}", app.FileHandler.HandleDeleteFile)
		r.Get("/ListFiles", app.FileHandler.HandleListFiles)
	})
	return router
}
