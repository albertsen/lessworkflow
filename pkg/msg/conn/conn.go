package conn

import (
	"os"

	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
)

func Connect() error {
	addr := os.Getenv("QUEUE_ADDR")
	if addr == "" {
		addr = "amqp://guest:guest@localhost:5672/"
	}
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	return err
}

func Channel(queueName string) (*amqp.Channel, error) {
	if channel == nil {
		var err error
		channel, err = conn.Channel()
		if err != nil {
			return nil, err
		}
	}
	_, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return channel, nil
}
