package msg

import (
	"encoding/json"
	"os"

	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
)

type Queue struct {
	Name    string
	Channel *amqp.Channel
}

type Message struct {
	Err     error
	Content interface{}
}

func Connect() error {
	addr := os.Getenv("QUEUE_ADDR")
	if addr == "" {
		addr = "amqp://guest:guest@localhost:5672/"
	}
	var err error
	conn, err = amqp.Dial(addr)
	return err
}

func Close() error {
	conn = nil
	return conn.Close()
}

func (Q *Queue) Publish(Message interface{}) error {
	data, err := json.Marshal(Q)
	if err != nil {
		return err
	}
	return Q.Channel.Publish(
		"",     // exchange
		Q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
}

func (Q *Queue) Consume(NewContent func() interface{}, quit chan bool) (chan Message, error) {
	msgs, err := Q.Channel.Consume(
		Q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}
	msgChan := make(chan Message)
	go func() {
		for {
			select {
			case d := <-msgs:
				content := NewContent()
				err := json.Unmarshal(d.Body, &content)
				if err != nil {
					msgChan <- Message{Err: err}
				} else {
					msgChan <- Message{Content: &content}
				}
			case <-quit:
				return
			}
		}
	}()
	return msgChan, nil
}

func OpenQueue(Name string) (*Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = channel.QueueDeclare(
		Name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // argument
	)
	if err != nil {
		return nil, err
	}
	return &Queue{Name: Name, Channel: channel}, nil
}
