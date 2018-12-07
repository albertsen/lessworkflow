package msg

import (
	"log"

	"github.com/golang/protobuf/jsonpb"
)

func (con *Connection) PublishProtobuf(Topic string, Message interface{}) error {
	json, err := jsonpb.Marshal(Message)
	if err != nil {
		log.Printf("Error marshalling message [%s]: %v", err, Message)
		return err
	}
	return con.PublishBytes(Topic, json)
}

func (con *Connection) PublishBytes(Topic string, Message []byte) error {
	log.Printf("Publishing message to topic [%s]: %s", Topic, string(Message))
	if err := con.NatsConn.Publish(Topic, Message); err != nil {
		log.Printf("Error publishing message [%s]: %s", err, string(Message))
	}
	return nil
}
