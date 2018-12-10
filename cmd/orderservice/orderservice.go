package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/albertsen/lessworkflow/pkg/msg"

	pbAction "github.com/albertsen/lessworkflow/gen/proto/action"
	pbOrder "github.com/albertsen/lessworkflow/gen/proto/order"
	proto "github.com/golang/protobuf/proto"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	serviceAddr = flag.String("m", "nats://localhost:4222", "Address of messaging server.")
	listenAddr  = flag.String("a", ":50051", "Address to listen on.")
	topic       = flag.String("t", "actions", "Message topic")
	help        = flag.Bool("h", false, "This message.")
)

type service struct {
	conn *msg.Connection
}

func (service *service) CreateOrder(ctx context.Context, order *pbOrder.Order) (*pbOrder.OrderCreatedResponse, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	order.Id = newUUID.String()
	newUUID, err = uuid.NewV4()
	if err != nil {
		return nil, err
	}
	processID := newUUID.String()
	orderData, err := proto.Marshal(order)
	if err != nil {
		return nil, err
	}
	actionRequest := pbAction.Request{
		Name:      "createOrder",
		ProcessId: processID,
		Payload: &pbAction.Payload{
			Id:   order.Id,
			Type: "order",
			Content: &any.Any{
				TypeUrl: proto.MessageName(order),
				Value:   orderData,
			},
		},
	}
	service.conn.PublishProtobuf(*topic, &actionRequest)
	log.Printf("Order created with ID: " + order.Id)
	return &pbOrder.OrderCreatedResponse{OrderId: order.Id, ProcessId: processID}, nil
}

func main() {

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	conn, err := msg.Connect(*serviceAddr)
	if err != nil {
		log.Fatalf("Failed to connect to messaging server: %s", err)
	}

	pbOrder.RegisterOrderServiceServer(s, &service{conn})

	reflection.Register(s)
	log.Printf("Listening on address: %s", listenAddr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
