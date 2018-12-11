package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	pb "github.com/albertsen/lessworkflow/gen/proto/order"

	ds "cloud.google.com/go/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	listenAddr = flag.String("a", ":50051", "Address to listen on.")
	help       = flag.Bool("h", false, "This message.")
)

type service struct {
	DSClient *ds.Client
}

func (service *service) SaveOrder(ctx context.Context, order *pb.Order) (*pb.SaveOrderResponse, error) {
	key := ds.NameKey("order", order.Id, nil)
	_, err := service.DSClient.Put(ctx, key, order)
	return &pb.SaveOrderResponse{OrderId: order.Id, Created: false, Order: order}, err
}

func (service *service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	key := ds.NameKey("order", req.OrderId, nil)
	var order pb.Order
	err := service.DSClient.Get(ctx, key, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (service *service) DeleteOrder(context.Context, *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	return nil, nil
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
	server := grpc.NewServer()

	ctx := context.Background()
	dsClient, err := ds.NewClient(ctx, "sap-se-commerce-arch")
	if err != nil {
		log.Fatalf("Failed to create new Cloud Datastore client: %s", err)
	}
	pb.RegisterOrderStorageServiceServer(server, &service{DSClient: dsClient})

	reflection.Register(server)
	log.Printf("Listening on address: %s", *listenAddr)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
