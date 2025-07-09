package exceptions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dsniels/storage-service/internal/utils"
)

var ErrorNotFound = errors.New("404: Data Not Fount")
var ErrorBadRequest = errors.New("400: Bad Request")

func HandleError(w http.ResponseWriter, err interface{}) {

	str := fmt.Sprint(err)
	arr := strings.Split(str, ":")

	key := arr[0]
	message := strings.Trim(arr[1], " ")

	switch key {
	case "400":
		utils.WriteResponse(w, http.StatusBadRequest, message)
	case "404":
		utils.WriteResponse(w, http.StatusNotFound, message)
	default:
		utils.WriteResponse(w, http.StatusInternalServerError, message)
	}

}

func ThrowException(statusCode int, message string) {
	err := errors.New(message)
	err = fmt.Errorf("%v: %w", statusCode, err)
	panic(err)
}
