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
	uuid "github.com/satori/go.uuid"

	ds "cloud.google.com/go/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func (service *service) CreateOrder(ctx context.Context, req *oss.CreateOrderRequest) (*oss.CreateOrderResponse, error) {
	order := req.Order
	if order.Id != "" {
		return &oss.CreateOrderResponse{}, status.New(codes.InvalidArgument, "New order cannot have an ID").Err()
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return &oss.CreateOrderResponse{}, nil
	}
	order.Id = uuid.String()
	key := ds.NameKey("order", order.Id, nil)
	mut := ds.NewInsert(key, order)
	_, err = service.DSClient.Mutate(ctx, mut)
	if err != nil {
		return &oss.CreateOrderResponse{}, err
	}
	return &oss.CreateOrderResponse{OrderId: order.Id}, nil
}

func (service *service) UpdateOrder(ctx context.Context, req *oss.UpdateOrderRequest) (*oss.UpdateOrderResponse, error) {
	order := req.Order
	if order.Id == "" {
		return &oss.UpdateOrderResponse{}, status.New(codes.InvalidArgument, "Order is missing ID").Err()
	}
	key := ds.NameKey("order", order.Id, nil)
	mut := ds.NewUpdate(key, order)
	_, err := service.DSClient.Mutate(ctx, mut)
	if err != nil {
		return &oss.UpdateOrderResponse{}, err
	}
	return &oss.UpdateOrderResponse{}, nil
}

func (service *service) GetOrder(ctx context.Context, req *oss.GetOrderRequest) (*oss.GetOrderResponse, error) {
	if req.OrderId == "" {
		return &oss.GetOrderResponse{}, status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
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
	if req.OrderId == "" {
		return &oss.DeleteOrderResponse{}, status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
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
