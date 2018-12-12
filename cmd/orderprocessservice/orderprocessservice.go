package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	ops "github.com/albertsen/lessworkflow/gen/proto/orderprocessservice"
	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	uuid "github.com/satori/go.uuid"

	ds "cloud.google.com/go/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port                    = flag.Int("p", 50051, "Port to listen on.")
	help                    = flag.Bool("h", false, "This message.")
	gcpProjectID            = os.Getenv("LW_GCP_PROJECT_ID")
	orderStorageServiceAddr = os.Getenv("LW_ORDER_STORAGE_SERVICE")
	listenAddr              string
)

type service struct {
	OrderStorageService oss.OrderStorageServiceServer
}

func (service *service) PlaceOrder(ctx context.Context, req *ops.PlaceOrderRequest) (*ops.PlaceOrderResponse, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return &ops.PlaceOrderResponse{}, err
	}
	orderID := uuid.String()
	order := req.Order
	order.Id = orderID
	_, err = service.OrderStorageService.SaveOrder(ctx, &oss.SaveOrderRequest{Order: order})
	if err != nil {
		return &ops.PlaceOrderResponse{}, err
	}
	return &ops.PlaceOrderResponse{OrderId: orderID}, nil
}

func (service *service) GetOrder(ctx context.Context, req *ops.GetOrderRequest) (*ops.GetOrderResponse, error) {
	res, err := service.OrderStorageService.GetOrder(ctx, &oss.GetOrderRequest{OrderId: req.OrderId})
	if err != nil {
		return &ops.GetOrderResponse{}, err
	}
	return &ops.GetOrderResponse{Order: res.Order}, nil
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
