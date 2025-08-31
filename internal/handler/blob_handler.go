package handler

import (
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
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

	w.Header().Add("Cache-Control", "public, max-age=31536000")
	w.Header().Add("Content-Type", "video/mp4")
	http.ServeContent(w, r, string(int32(id)), time.Time{}, streamer)

}

func (c *BlobHandler) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024*1024)
	contentType := r.Header.Get("Content-Type")
	_, params, _ := mime.ParseMediaType(contentType)
	m := multipart.NewReader(r.Body, params["boundary"])

	for {
		part, err := m.NextPart()
		if err == io.EOF {
			break
		}
		url, err := c.store.UploadBlob(r.Context(), part.FileName(), part, part.Header.Values("content-type")[0])
		if err != nil {
			exceptions.ThrowInternalServerError("uploading to azure", err.Error())
		}
		utils.WriteResponse(w, http.StatusOK, utils.Response{"data": url})
	}
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
