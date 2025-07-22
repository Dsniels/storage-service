package controllers

import (
	"log"
	"net/http"

	"time"

	exceptions "github.com/dsniels/storage-service/internal/Exceptions"
	"github.com/dsniels/storage-service/internal/params"
	store "github.com/dsniels/storage-service/internal/storage"
	"github.com/dsniels/storage-service/internal/utils"
	pb "github.com/dsniels/storage-service/proto"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	store     store.IStore
	stream    store.IStream
	rpcClient pb.CursosProtoServiceClient
}

type IController interface {
	HandleUploadFile(w http.ResponseWriter, r *http.Request)
	HandleStreamFile(w http.ResponseWriter, r *http.Request)
	HandleListFiles(w http.ResponseWriter, r *http.Request)
	HandleDeleteFile(w http.ResponseWriter, r *http.Request)
}

func (c *Controller) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 3*1024*1024*1024)
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
	var b = struct{ UserId string }{}
	utils.GetParamsFromUrl(r, b)
	response, err := c.rpcClient.CheckUserAccess(r.Context(), &pb.CursoAccessRequest{UserId: b.UserId})
	if err != nil {
		exceptions.ThrowInternalServerError(err.Error())
	}

	if response.Ok {
		streamer, err := c.stream.GetStream(r.Context(), id)
		if err != nil {
			exceptions.ThrowInternalServerError(err.Error())
		}
		http.ServeContent(w, r, id, time.Time{}, streamer)
	} else {
		exceptions.ThrowBadRequest("you cannot see this")
	}

}

func (c *Controller) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	response, err := c.rpcClient.SayHi(r.Context(), &pb.HiRequest{Name: "Daniel"})
	if err != nil {
		log.Panicln(err)
	}
	log.Println(response.Message)
	var p params.FileList
	err = utils.GetParamsFromUrl(r, &p)
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

func NewController(store store.IStore, stream store.IStream, rpc pb.CursosProtoServiceClient) *Controller {
	return &Controller{
		store:     store,
		rpcClient: rpc,
		stream:    stream,
	}
}
