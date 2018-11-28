package msg

import (
	"encoding/json"
	"log"
	"time"

	nats "github.com/nats-io/go-nats"
)

type Subscripton struct {
	Subscription *nats.Subscription
}

func (c *Connection) Subscribe(Topic string) *Subscripton {
	log.Printf("Subscribing to topic: %s", Topic)
	sub, err := c.NatsConn.SubscribeSync(Topic)
	if err != nil {
		log.Fatal(err)
	}
	return &Subscripton{Subscription: sub}
}

func (s *Subscripton) NextMessage(Message interface{}) bool {
	msg, err := s.Subscription.NextMsg(1 * time.Hour)
	if err != nil {
		if err == nats.ErrTimeout {
			return false
		}
		log.Fatal(err)
	}
	log.Printf("Raw message received: %s", msg.Data)
	json.Unmarshal([]byte(msg.Data), &Message)
	log.Printf("Unmarshalled message: %v", Message)
	return true
}
