package controllers

import (
	"log"
	"net/http"

	"time"

	exceptions "github.com/dsniels/storage-service/internal/Exceptions"
	"github.com/dsniels/storage-service/internal/params"
	"github.com/dsniels/storage-service/internal/storage"
	"github.com/dsniels/storage-service/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	store  storage.IStore
	stream storage.IStream
}

type IController interface {
	HandleUploadFile(w http.ResponseWriter, r *http.Request)
	HandleStreamFile(w http.ResponseWriter, r *http.Request)
	HandleListFiles(w http.ResponseWriter, r *http.Request)
	HandleDeleteFile(w http.ResponseWriter, r *http.Request)
}

func (c *Controller) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 3*1024*1024*1024)
	log.Println(r.Method)
	file, file_header, err := r.FormFile("file")
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}
	buff := make([]byte, file_header.Size)
	_, err = file.Read(buff)
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}

	url, err := c.store.UploadFile(r.Context(), file_header.Filename, buff, file_header.Header.Get("content-type"))
	if err != nil {
		exceptions.ThrowInternalServerError(err.Error())
	}
	utils.WriteResponse(w, http.StatusOK, utils.Response{"data": url})
}

func (c *Controller) HandleStreamFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		exceptions.ThrowBadRequest("Must pass an id")
	}

	streamer, err := c.stream.GetStream(r.Context(), id)
	if err != nil {
		exceptions.ThrowInternalServerError(err.Error())
	}
	http.ServeContent(w, r, id, time.Time{}, streamer)

}

func (c *Controller) HandleListFiles(w http.ResponseWriter, r *http.Request) {

	var p params.FileList
	err := utils.GetParamsFromUrl(r, &p)
	if err != nil {
		exceptions.ThrowInternalServerError(err.Error())
	}

	files, err := c.store.GetFiles(r.Context(), "temp", p.Prefix)
	if err != nil {
		exceptions.ThrowInternalServerError(err.Error())
	}
	data := struct {
		Files  *[]string
		Params params.FileList
	}{Files: files, Params: p}
	utils.WriteResponse(w, http.StatusOK, utils.Response{"data": data})

}

func (c *Controller) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdFromUrl(r)
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}

	err = c.store.DeleteFile(r.Context(), id, "")

	if err != nil {
		exceptions.ThrowInternalServerError()
	}

	utils.WriteResponse(w, http.StatusOK, &utils.Response{"data": struct{}{}})

}

func NewController(store storage.IStore, stream storage.IStream) *Controller {
	return &Controller{
		store:  store,
		stream: stream,
	}
}
