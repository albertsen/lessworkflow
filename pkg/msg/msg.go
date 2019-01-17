package msg

import (
	"encoding/json"
	"log"
	"os"

	"github.com/assembla/cony"
	"github.com/streadway/amqp"
)

var (
	client *cony.Client
)

type Publisher struct {
	Pub *cony.Publisher
}

func (P *Publisher) Publish(Message interface{}) error {
	data, err := json.Marshal(Message)
	if err != nil {
		return err
	}
	return P.Pub.Publish(amqp.Publishing{
		Body: data,
	})
}

func init() {
	addr := os.Getenv("QUEUE_ADDR")
	if addr == "" {
		addr = "amqp://guest:guest@localhost:5672/"
	}
	client = cony.NewClient(
		cony.URL(addr),
		cony.Backoff(cony.DefaultBackoff),
	)
}

func NewPublisher(Name string) *Publisher {
	exc := cony.Exchange{
		Name:       Name,
		Kind:       "fanout",
		AutoDelete: false,
		Durable:    true,
	}
	client.Declare([]cony.Declaration{
		cony.DeclareExchange(exc),
	})
	publisher := cony.NewPublisher(exc.Name, "")
	client.Publish(publisher)
	return &Publisher{Pub: publisher}
}

func StartLoop() {
	go func() {
		for client.Loop() {
			select {
			case err := <-client.Errors():
				log.Printf("Reconnecting to RabbitMQ after error: %s", err)
			}
		}
	}()
}

func Close() {
	client.Close()
}
