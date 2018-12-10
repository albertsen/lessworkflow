package msg

import (
	"log"

	"github.com/golang/protobuf/proto"
)

func (con *Connection) PublishProtobuf(Topic string, Message proto.Message) error {
	data, err := proto.Marshal(Message)
	if err != nil {
		return err
	}
	return con.PublishBytes(Topic, data)
}

func (con *Connection) PublishBytes(Topic string, Message []byte) error {
	log.Printf("Publishing message to topic [%s]: %s", Topic, string(Message))
	return con.NatsConn.Publish(Topic, Message)
}
