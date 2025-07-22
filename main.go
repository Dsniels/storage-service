package main

import (
	"context"
	"log"
	"net/http"
	"os"

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
	defer services.Conn.Close()
	go services.Queue.DeleteFileConsumer(context.Background())

	router := router.InitRoutes(services)
	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT_STORE"),
		Handler: router,
	}
	log.Println("Running server....")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
