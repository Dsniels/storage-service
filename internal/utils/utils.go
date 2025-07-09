package utils

import (
	"encoding/json"
	"net/http"
)

type Response map[string]any

func WriteResponse(w http.ResponseWriter, status int, data any) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json")

	w.WriteHeader(status)
	w.Write(json)
	return nil
}
