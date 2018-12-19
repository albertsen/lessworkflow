package main

import (
	"context"
	"log"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	"github.com/albertsen/lessworkflow/pkg/grpc/grpcserver"
	"github.com/golang/protobuf/ptypes"

	uuid "github.com/satori/go.uuid"

	"github.com/go-pg/pg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	DB *pg.DB
}

func (service *service) CreateOrder(ctx context.Context, req *oss.CreateOrderRequest) (*oss.CreateOrderResponse, error) {
	order := req.Order
	if req.Order == nil {
		return &oss.CreateOrderResponse{}, status.New(codes.InvalidArgument, "No order provided").Err()
	}
	if order.Id != "" {
		return &oss.CreateOrderResponse{}, status.New(codes.InvalidArgument, "New order cannot have an ID").Err()
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return &oss.CreateOrderResponse{}, nil
	}
	order.Id = uuid.String()
	now := ptypes.TimestampNow()
	order.TimeCreated = now
	order.TimeUpdated = now
	order.Status = "CREATED"
	order.Version = 1
	if err = service.DB.Insert(order); err != nil {
		log.Printf("Error creating order: %s", err)
		return &oss.CreateOrderResponse{}, err
	}
	return &oss.CreateOrderResponse{OrderId: order.Id}, nil
}

func (service *service) UpdateOrder(ctx context.Context, req *oss.UpdateOrderRequest) (*oss.UpdateOrderResponse, error) {
	order := req.Order
	if order.Id == "" {
		return &oss.UpdateOrderResponse{}, status.New(codes.InvalidArgument, "Order is missing ID").Err()
	}
	order.TimeUpdated = ptypes.TimestampNow()
	if err := service.DB.Update(order); err != nil {
		return &oss.UpdateOrderResponse{}, err
	}
	return &oss.UpdateOrderResponse{}, nil
}

func (service *service) GetOrder(ctx context.Context, req *oss.GetOrderRequest) (*oss.GetOrderResponse, error) {
	if req.OrderId == "" {
		return &oss.GetOrderResponse{}, status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
	var order = od.Order{
		Id: req.OrderId,
	}
	if err := service.DB.Select(order); err != nil {
		return &oss.GetOrderResponse{}, err
	}
	return &oss.GetOrderResponse{Order: &order}, nil
}

func (service *service) DeleteOrder(ctx context.Context, req *oss.DeleteOrderRequest) (*oss.DeleteOrderResponse, error) {
	if req.OrderId == "" {
		return &oss.DeleteOrderResponse{}, status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
	var order = od.Order{
		Id: req.OrderId,
	}
	if err := service.DB.Delete(order); err != nil {
		return &oss.DeleteOrderResponse{}, err
	}
	return &oss.DeleteOrderResponse{}, nil
}

func main() {
	grpcserver.StartServer(func(server *grpc.Server) {
		db := pg.Connect(&pg.Options{
			User:     "lwadmin",
			Database: "lessworkflow",
		})
		// defer db.Close()
		oss.RegisterOrderStorageServiceServer(server, &service{DB: db})
	})

}
