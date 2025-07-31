package app

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

func getBlobClient() *azblob.Client {
	connection := os.Getenv("connection")
	client, err := azblob.NewClientFromConnectionString(connection, &azblob.ClientOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getFileClient() *share.Client {
	connection := os.Getenv("connection")
	client, err := share.NewClientFromConnectionString(connection, "temp", nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime)
}
