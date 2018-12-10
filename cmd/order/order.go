package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pbOrder "github.com/albertsen/lessworkflow/gen/proto/order"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
)

var (
	serviceAddr = flag.String("s", "localhost:50051", "Address of Order Service.")
	help        = flag.Bool("h", false, "This message.")
)

func main() {
	flag.Parse()
	command := flag.Arg(0)
	orderFileName := flag.Arg(1)
	if *help || command != "create" || orderFileName == "" {
		fmt.Printf("Usage: %s [options] create <order file>\n\nValid options:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var order pbOrder.Order
	orderFile, err := os.Open(orderFileName)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer orderFile.Close()
	err = jsonpb.Unmarshal(orderFile, &order)
	if err != nil {
		log.Fatalf("Error unmarshaling file: %s", err)
	}

	conn, err := grpc.Dial(*serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer conn.Close()
	client := pbOrder.NewOrderServiceClient(conn)

	r, err := client.CreateOrder(context.Background(), &order)
	if err != nil {
		log.Fatalf("Error placing order: %s", err)
	}
	log.Printf("Order created with ID: %s", r.OrderId)
}
