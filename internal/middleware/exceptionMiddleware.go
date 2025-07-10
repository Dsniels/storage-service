package middleware

import (
	"log"
	"net/http"

	exceptions "github.com/dsniels/storage-service/internal/Exceptions"
)

func Exception(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				exceptions.HandleError(w, err)
			}
		}()

		next.ServeHTTP(w, r)
	})

}
