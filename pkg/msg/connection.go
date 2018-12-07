package msg

import (
	"log"

	nats "github.com/nats-io/go-nats"
)

type Connection struct {
	NatsConn *nats.Conn
}

func (con *Connection) Close() {
	log.Printf("Closing connection to NATS server")
	con.NatsConn.Flush()
	con.NatsConn.Close()
}

func Connect(URL string) (*Connection, error) {
	log.Printf("Connecting to NATS server at: %s", URL)
	nc, err := nats.Connect(URL)
	if err != nil {
		return nil, err
	}
	return &Connection{NatsConn: nc}, nil
}
