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

func (p *Publisher) Publish(content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return p.Pub.Publish(amqp.Publishing{
		Body:            data,
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
	})
}

type Consumer struct {
	Cns *cony.Consumer
}

func (c *Consumer) Consume(newContent func() interface{}, processContent func(interface{}) error) {
	for client.Loop() {
		select {
		case msg := <-c.Cns.Deliveries():
			content := newContent()
			err := json.Unmarshal(msg.Body, content)
			if err != nil {
				log.Printf("Error unmarshaling message content: %s", err)
				break
			}
			err = processContent(content)
			if err != nil {
				log.Printf("Error prcessing content of message: %s", err)
				// TODO: Put in error queue
			}
		case err := <-c.Cns.Errors():
			log.Printf("Consumer error: %v\n", err)
		case err := <-client.Errors():
			log.Printf("Client error: %v\n", err)
		}
	}
}

func init() {
	addr := os.Getenv("MSG_SERVER_URL")
	if addr == "" {
		addr = "amqp://guest:guest@localhost:5672/"
	}
	client = cony.NewClient(
		cony.URL(addr),
		cony.Backoff(cony.DefaultBackoff),
	)
}

func NewPublisher(name string) *Publisher {
	exc := cony.Exchange{
		Name:       name,
		Kind:       "direct",
		Durable:    true,
		AutoDelete: false,
	}
	client.Declare([]cony.Declaration{
		cony.DeclareExchange(exc),
	})
	publisher := cony.NewPublisher(exc.Name, "")
	client.Publish(publisher)
	return &Publisher{Pub: publisher}
}

func NewConsumer(name string) *Consumer {
	que := &cony.Queue{
		AutoDelete: false,
		Durable:    true,
	}
	exc := cony.Exchange{
		Name:       name,
		Kind:       "direct",
		AutoDelete: true,
		Durable:    true,
	}
	bnd := cony.Binding{
		Queue:    que,
		Exchange: exc,
		Key:      "",
	}
	client.Declare([]cony.Declaration{
		cony.DeclareQueue(que),
		cony.DeclareExchange(exc),
		cony.DeclareBinding(bnd),
	})
	cns := cony.NewConsumer(
		que,
	)
	return &Consumer{Cns: cns}
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
