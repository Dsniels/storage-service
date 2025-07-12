package app

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func getAzureClient() *azblob.Client {

	connection := os.Getenv("connection")
	client, err := azblob.NewClientFromConnectionString(connection, &azblob.ClientOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime)
}
