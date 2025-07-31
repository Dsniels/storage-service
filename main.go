package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dsniels/storage-service/internal/app"
	"github.com/dsniels/storage-service/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	services := app.InitServices()


	router := router.InitRoutes(services)
	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT_STORE"),
		Handler: router,
	}
	log.Println("Running server....")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
