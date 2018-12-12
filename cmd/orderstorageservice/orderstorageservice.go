package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"

	ds "cloud.google.com/go/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port         = flag.Int("p", 50051, "Port to listen on.")
	help         = flag.Bool("h", false, "This message.")
	gcpProjectID = os.Getenv("LW_GCP_PROJECT_ID")
	listenAddr   string
)

type service struct {
	DSClient *ds.Client
}

func (service *service) SaveOrder(ctx context.Context, req *oss.SaveOrderRequest) (*oss.SaveOrderResponse, error) {
	order := req.Order
	key := ds.NameKey("order", order.Id, nil)
	_, err := service.DSClient.Put(ctx, key, order)
	return &oss.SaveOrderResponse{Order: order}, err
}

func (service *service) GetOrder(ctx context.Context, req *oss.GetOrderRequest) (*oss.GetOrderResponse, error) {
	key := ds.NameKey("order", req.OrderId, nil)
	var order od.Order
	err := service.DSClient.Get(ctx, key, &order)
	if err == ds.ErrNoSuchEntity {
		return &oss.GetOrderResponse{}, nil
	}
	if err != nil {
		return &oss.GetOrderResponse{}, err
	}
	return &oss.GetOrderResponse{Order: &order}, nil
}

func (service *service) DeleteOrder(ctx context.Context, req *oss.DeleteOrderRequest) (*oss.DeleteOrderResponse, error) {
	key := ds.NameKey("order", req.OrderId, nil)
	service.DSClient.Delete(ctx, key)
	return &oss.DeleteOrderResponse{}, nil
}

func init() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}
	listenAddr = ":" + strconv.Itoa(*port)
	if gcpProjectID == "" {
		gcpProjectID = "sap-se-commerce-arch"
	}
}

func main() {

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()

	ctx := context.Background()
	dsClient, err := ds.NewClient(ctx, gcpProjectID)
	if err != nil {
		log.Fatalf("Failed to create new Cloud Datastore client: %s", err)
	}
	oss.RegisterOrderStorageServiceServer(server, &service{DSClient: dsClient})

	reflection.Register(server)
	log.Printf("Listening on address: %s", listenAddr)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
