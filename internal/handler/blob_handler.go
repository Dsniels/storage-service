package handler

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	exceptions "github.com/dsniels/storage-service/internal/Exceptions"
	"github.com/dsniels/storage-service/internal/params"
	store "github.com/dsniels/storage-service/internal/storage"
	"github.com/dsniels/storage-service/internal/utils"
	pb "github.com/dsniels/storage-service/proto"
	"github.com/go-chi/chi/v5"
)

type BlobHandler struct {
	store     store.IStore
	stream    store.IStream
	rpcClient pb.CursosProtoServiceClient
	cache     sync.Map
}

func (c *BlobHandler) HandleStreamFile(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdFromUrl(r)
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}
	response, err := c.rpcClient.GetCursoByID(r.Context(), &pb.GetCursoRequest{Id: int32(id)})
	if err != nil {
		slog.Error("Error GetCurso grpc: ", err)
	}
	video := response.GetVideo()
	videoId, err := c.store.GetFileIdFromURL(r.Context(), video)
	if err != nil {
		slog.Error("Getting video id: ", err)
	}
	streamer, err := c.stream.GetStream(r.Context(), *videoId)
	if err != nil {
		slog.Error("Getting Stream: ", err)
		exceptions.ThrowInternalServerError(err.Error())
	}
	w.Header().Add("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, string(int32(id)), time.Time{}, streamer)

}

func (c *BlobHandler) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
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

func (c *BlobHandler) HandleListFiles(w http.ResponseWriter, r *http.Request) {
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

func (c *BlobHandler) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		exceptions.ThrowBadRequest("Invalid ID")
	}

	if err := c.store.DeleteFile(r.Context(), id, ""); err != nil {
		exceptions.ThrowInternalServerError(err.Error())
	}

	utils.WriteResponse(w, http.StatusOK, &utils.Response{"data": struct{}{}})
}

func NewBlobHandler(store store.IStore, stream store.IStream, rpc pb.CursosProtoServiceClient) *BlobHandler {

	return &BlobHandler{
		store:     store,
		rpcClient: rpc,
		stream:    stream,
	}
}
