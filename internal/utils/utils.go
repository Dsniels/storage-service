package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
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
func GetIdFromUrl(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return 0, fmt.Errorf("must provide an id")
	}

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatalln("error parsing int", err)
		return 0, fmt.Errorf("must provide an id")
	}

	return i, nil
}

func GetQueryFromUrl(r *http.Request) (string, error) {

	return "", nil
}

func GetParamsFromUrl(r *http.Request, s interface{}) error {
	val := reflect.ValueOf(s)
	log.Println(val)
	if val.Kind() != reflect.Pointer && val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("must pass a pointer to an struct")
	}

	el := val.Elem()
	t := el.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := r.URL.Query().Get(field.Name)
		if el.Field(i).CanSet() {
			el.Field(i).SetString(val)
		}
	}

	return nil

}
