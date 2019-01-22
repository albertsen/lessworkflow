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
	log.Printf("Publishing message: %s", string(data))
	return p.Pub.Publish(amqp.Publishing{
		Body:            data,
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
	})
}

type Consumer struct {
	Cns *cony.Consumer
}

func (c *Consumer) Consume(newContentStruct func() interface{}, processContent func(interface{}) error, done chan bool) {
	for client.Loop() {
		log.Println("Starting consumer loop")
		select {
		case msg := <-c.Cns.Deliveries():
			log.Printf("Received message: %s", string(msg.Body))
			content := newContentStruct()
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
			msg.Ack(false)
		case err := <-c.Cns.Errors():
			log.Printf("Consumer error: %v\n", err)
		case err := <-client.Errors():
			log.Printf("Client error: %v\n", err)
		case <-done:
			return
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

func declareQueue(name string) *cony.Queue {
	que := &cony.Queue{
		Name:       name,
		AutoDelete: false,
		Durable:    true,
	}
	exc := cony.Exchange{
		Name:       name,
		Kind:       "fanout",
		AutoDelete: false,
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
	return que
}

func NewPublisher(name string) *Publisher {
	declareQueue(name)
	publisher := cony.NewPublisher(name, "")
	client.Publish(publisher)
	return &Publisher{Pub: publisher}
}

func NewConsumer(name string) *Consumer {
	que := declareQueue(name)
	cns := cony.NewConsumer(que)
	return &Consumer{Cns: cns}
}

func StartConnectionLoop() {
	log.Printf("Starting connection loop")
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
