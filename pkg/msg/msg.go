package msg

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/streadway/amqp"
)

var (
	addr       string
	conn       *amqp.Connection
	chann      *amqp.Channel
	connMutex  = &sync.Mutex{}
	channMutex = &sync.Mutex{}
)

func Publish(queueName string, content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	log.Printf("Publishing message: %s", string(data))
	ch, err := channel()
	if err != nil {
		return err
	}
	return ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
		})
}

func Consume(queueName string,
	newContentStruct func() interface{},
	processContent func(interface{}) error,
	errChan chan error,
	doneChan chan bool) error {
	ch, err := channel()
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	select {
	case msg := <-msgs:
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
			errChan <- err
		}
		msg.Ack(false)
	case <-doneChan:
		return nil
	}
	return nil
}

func init() {
	addr := os.Getenv("MSG_SERVER_URL")
	if addr == "" {
		addr = "amqp://guest:guest@localhost:5672/"
	}
}

func declareQueue(name string) (*amqp.Queue, error) {
	ch, err := channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func closeChannel() {
	channMutex.Lock()
	defer channMutex.Unlock()
	if chann != nil {
		chann.Close()
		chann = nil
	}
}

func closeConnection() {
	connMutex.Lock()
	defer connMutex.Unlock()
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

func connection() (*amqp.Connection, error) {
	connMutex.Lock()
	defer connMutex.Unlock()
	if conn == nil {
		var err error
		conn, err = amqp.Dial(addr)
		if err != nil {
			conn = nil
			return nil, err
		}
	}
	return conn, nil
}

func channel() (*amqp.Channel, error) {
	channMutex.Lock()
	defer channMutex.Unlock()
	if chann == nil {
		con, err := connection()
		if err != nil {
			return nil, err
		}
		chann, err = con.Channel()
		if err != nil {
			return nil, err
		}
	}
	return chann, nil
}

func init() {
	addr = os.Getenv("MSG_SERVER_URL")
	if addr == "" {
		addr = "amqp://guest:guest@localhost:5672/"
	}
}

func Close() {
	closeChannel()
	closeConnection()
}
