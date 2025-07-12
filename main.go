package main

import (
	"log"
	"net/http"

	"github.com/dsniels/storage-service/internal/app"
	"github.com/dsniels/storage-service/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	services := app.InitServices()
	router := router.InitRoutes(services)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("Running server....")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
