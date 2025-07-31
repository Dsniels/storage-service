package handler

import (
	"fmt"
	"net/http"

	exceptions "github.com/dsniels/storage-service/internal/Exceptions"
	store "github.com/dsniels/storage-service/internal/storage"
	"github.com/dsniels/storage-service/internal/utils"
)

type FileHandler struct {
	store store.IStore
}

func (f *FileHandler) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	file, file_header, err := r.FormFile("file")
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}
	buf := make([]byte, file_header.Size)

	_, err = file.Read(buf)
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}

	url, err := f.store.UploadFile(r.Context(), file_header.Filename, buf, file_header.Header.Get("content-type"))
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}

	fmt.Println(url)
	utils.WriteResponse(w, http.StatusOK, utils.Response{"data": url})

}

func (f *FileHandler) HandleListFiles(w http.ResponseWriter, r *http.Request) {

	files, err := f.store.GetFiles(r.Context(), "files", "")
	if err != nil {
		exceptions.ThrowBadRequest(err.Error())
	}
	utils.WriteResponse(w, http.StatusOK, utils.Response{"data": files}) }


func (f *FileHandler) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {

}

func NewFileHandler(store store.IStore) *FileHandler {

	return &FileHandler{
		store: store,
	}
}
