package handler

import "net/http"

type IStoreHandler interface {
	HandleUploadFile(w http.ResponseWriter, r *http.Request)
	HandleListFiles(w http.ResponseWriter, r *http.Request)
	HandleDeleteFile(w http.ResponseWriter, r *http.Request)
}

type IFileHandler interface {
	IStoreHandler
}

type IBlobHandler interface {
	IStoreHandler
	IStreamHandler
}

type IStreamHandler interface { 
	HandleStreamFile(w http.ResponseWriter, r *http.Request) }
