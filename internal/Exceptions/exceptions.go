package exceptions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dsniels/storage-service/internal/utils"
)

var ErrorNotFound = errors.New("data not found")
var ErrorInternal = errors.New("something went wrong")
var ErrorBadRequest = errors.New("bad request")

func HandleError(w http.ResponseWriter, err interface{}) {

	str := fmt.Sprint(err)
	arr := strings.Split(str, ":")

	key := arr[0]
	val := arr[len(arr)-1]
	message := strings.Trim(val, " ")

	switch key {
	case "400":
		utils.WriteResponse(w, http.StatusBadRequest, &utils.Response{"message": message})
	case "404":
		utils.WriteResponse(w, http.StatusNotFound, utils.Response{"message": message})
	default:
		utils.WriteResponse(w, http.StatusInternalServerError, utils.Response{"message": message})
	}

}

func ThrowException(statusCode int, message string) {
	err := errors.New(message)
	err = fmt.Errorf("%v: %v", statusCode, err)
	panic(err)
}

func ThrowNotFound(messages ...string) {
	var err error = ErrorNotFound
	if len(messages) > 0 {
		err = errors.New(messages[0])
	}
	ThrowException(http.StatusNotFound, err.Error())

}

func ThrowBadRequest(messages ...string) {
	var err error = ErrorBadRequest
	if len(messages) > 0 {
		err = errors.New(messages[0])
	}
	ThrowException(http.StatusBadRequest, err.Error())

}
func ThrowInternalServerError(message ...string) {
	var err error = ErrorInternal
	if len(message) > 0 {
		err = errors.New(message[0])
	}
	ThrowException(http.StatusInternalServerError, err.Error())
}
