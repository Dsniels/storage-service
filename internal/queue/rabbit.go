package queue

import (
	"context"
	"encoding/json"
	"log"
	"os"

	store "github.com/dsniels/storage-service/internal/storage"
	"github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	conn        *amqp091.Connection
	ch          *amqp091.Channel
	blobStorage store.IStore
}

func (r *Rabbit) DeleteFileConsumer(ctx context.Context) {
	defer r.ch.Close()

	err := r.ch.ExchangeDeclare("DeleteBlob", amqp091.ExchangeFanout, false, true, false, false, nil)
	if err != nil {
		log.Fatalf("Exchange Deglare %v", err)
	}

	queue, err := r.ch.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Queue Deglare %v", err)
	}

	err = r.ch.QueueBind(queue.Name, "", "DeleteBlob", false, nil)
	if err != nil {
		log.Fatalf("Queue Binding %v", err)
	}

	msg, err := r.ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Consume %v", err)
	}

	var stop chan struct{}

	go func() {
		for d := range msg {
			message := new(DeleteBlob)
			json.Unmarshal(d.Body, message)
			id, _ := r.blobStorage.GetBlobIdFromURL(ctx, message.Url)
			r.blobStorage.DeleteFile(ctx, *id, "")
		}
	}()

	<-stop
}

func NewRabbit(svc store.IStore) (*Rabbit, error) {
	conn, err := amqp091.Dial(os.Getenv("RABBIT_CONN"))
	if err != nil {
		log.Fatalf("Error in connection %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error in channel %v", err)
	}
	return &Rabbit{
		conn:        conn,
		ch:          ch,
		blobStorage: svc,
	}, nil

}
