package router

import (
	"github.com/dsniels/storage-service/internal/app"
	"github.com/dsniels/storage-service/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(app *app.App) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Exception)

	router.Group(func(r chi.Router) {
		r.Post("/UploadFile", app.Controller.HandleUploadFile)
		r.Get("/Stream/{id}", app.Controller.HandleStreamFile)
		r.Get("/ListFiles", app.Controller.HandleListFiles)
	})

	return router
}
