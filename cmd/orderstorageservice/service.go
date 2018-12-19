package main

import (
	"context"
	"log"

	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	"github.com/albertsen/lessworkflow/pkg/grpc/grpcserver"
	"github.com/go-pg/pg"
	"github.com/golang/protobuf/ptypes"
	uuid "github.com/satori/go.uuid"
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
		log.Printf("ERROR generating UUID: %s", err)
		return &oss.CreateOrderResponse{}, err
	}
	order.Id = uuid.String()
	now := ptypes.TimestampNow()
	order.TimeCreated = now
	order.TimeUpdated = now
	if order.TimePlaced == nil {
		order.TimePlaced = now
	}
	order.Status = "CREATED"
	order.Version = 1
	orderDTO, err := newOrderDTO(order)
	if err != nil {
		log.Printf("ERROR converting protbuf Order into OrderDTO: %s", err)
		return &oss.CreateOrderResponse{}, err
	}
	if err = service.DB.Insert(orderDTO); err != nil {
		log.Printf("ERROR inserting Order into DB: %s", err)
		return &oss.CreateOrderResponse{}, err
	}
	return &oss.CreateOrderResponse{Order: order}, nil
}

func (service *service) UpdateOrder(ctx context.Context, req *oss.UpdateOrderRequest) (*oss.UpdateOrderResponse, error) {
	order := req.Order
	if order.Id == "" {
		return &oss.UpdateOrderResponse{}, status.New(codes.InvalidArgument, "Order is missing ID").Err()
	}
	order.TimeUpdated = ptypes.TimestampNow()
	orderDTO, err := newOrderDTO(order)
	if err != nil {
		log.Printf("ERROR converting protbuf Order into OrderDTO: %s", err)
		return &oss.UpdateOrderResponse{}, err
	}
	if err := service.DB.Update(orderDTO); err != nil {
		if err == pg.ErrNoRows {
			return &oss.UpdateOrderResponse{}, status.New(codes.NotFound, "Order not found").Err()
		}
		log.Printf("ERROR updating Order in DB: %s", err)
		return &oss.UpdateOrderResponse{}, err
	}
	return &oss.UpdateOrderResponse{Order: order}, nil
}

func (service *service) GetOrder(ctx context.Context, req *oss.GetOrderRequest) (*oss.GetOrderResponse, error) {
	if req.OrderId == "" {
		return &oss.GetOrderResponse{}, status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
	orderDTO := &orderDTO{
		ID: req.OrderId,
	}
	if err := service.DB.Select(orderDTO); err != nil {
		if err == pg.ErrNoRows {
			return &oss.GetOrderResponse{}, nil
		}
		log.Printf("ERROR selecting Order from DB: %s", err)
		return &oss.GetOrderResponse{}, err
	}
	order, err := newOrderProto(orderDTO)
	if err != nil {
		log.Printf("ERROR converting Order DTO to Order proto: %s", err)
		return &oss.GetOrderResponse{}, err
	}
	return &oss.GetOrderResponse{Order: order}, nil
}

func (service *service) DeleteOrder(ctx context.Context, req *oss.DeleteOrderRequest) (*oss.DeleteOrderResponse, error) {
	if req.OrderId == "" {
		return &oss.DeleteOrderResponse{}, status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
	var orderDTO = orderDTO{
		ID: req.OrderId,
	}
	if err := service.DB.Delete(&orderDTO); err != nil {
		log.Printf("ERROR deleting Order from DB: %s", err)
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
