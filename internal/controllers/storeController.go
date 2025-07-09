package controllers

import (
	"log"
	"net/http"
	"time"

	exceptions "github.com/dsniels/storage-service/internal/Exceptions"
	"github.com/dsniels/storage-service/internal/storage"
	"github.com/dsniels/storage-service/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	store storage.IStore
}

type IController interface {
	HandleUploadFile(w http.ResponseWriter, r *http.Request)
	HandleStreamFile(w http.ResponseWriter, r *http.Request)
}

func NewController(store storage.IStore) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 3*1024*1024*1024)
	log.Println(r.Method)
	file, file_header, err := r.FormFile("file")
	if err != nil {
		log.Println("error FormFile: ", err)
		exceptions.ThrowException(http.StatusBadRequest, err.Error())
	}
	buff := make([]byte, file_header.Size)
	_, err = file.Read(buff)
	if err != nil {
		log.Println("error reading: ", err)
		exceptions.ThrowException(http.StatusBadRequest, err.Error())
	}

	url, err := c.store.UploadFile(r.Context(), file_header.Filename, buff, file_header.Header.Get("content-type"))
	if err != nil {
		log.Println("error store uploading : ", err)

		exceptions.ThrowException(http.StatusBadRequest, err.Error())
	}
	utils.WriteResponse(w, http.StatusOK, utils.Response{"data": url})
}

func (c *Controller) HandleStreamFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		exceptions.ThrowException(http.StatusBadRequest, "Must pass an id")
	}

	stream, err := c.store.GetFile(r.Context(), id)
	if err != nil {
		exceptions.ThrowException(http.StatusInternalServerError, err.Error())
	}



	http.ServeContent(w, r, id, time.Time{}, &storage.BlobReader{Ctx: r.Context(),  })

}
