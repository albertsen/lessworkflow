package main

import (
	"context"

	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	dao "github.com/albertsen/lessworkflow/pkg/db/daos/orderdao"
	"github.com/albertsen/lessworkflow/pkg/grpc/grpcserver"
	"google.golang.org/grpc"
)

type service struct {
}

func (service *service) CreateOrder(ctx context.Context, req *oss.CreateOrderRequest) (*oss.CreateOrderResponse, error) {
	order, err := dao.CreateOrder(req.Order)
	if err != nil {
		return &oss.CreateOrderResponse{}, err
	}
	return &oss.CreateOrderResponse{Order: order}, nil
}

func (service *service) UpdateOrder(ctx context.Context, req *oss.UpdateOrderRequest) (*oss.UpdateOrderResponse, error) {
	order, err := dao.UpdateOrder(req.Order)
	if err != nil {
		return &oss.UpdateOrderResponse{}, err
	}
	return &oss.UpdateOrderResponse{Order: order}, nil
}

func (service *service) GetOrder(ctx context.Context, req *oss.GetOrderRequest) (*oss.GetOrderResponse, error) {
	order, err := dao.GetOrder(req.OrderId)
	if err != nil {
		return &oss.GetOrderResponse{}, err
	}
	return &oss.GetOrderResponse{Order: order}, nil
}

func (service *service) DeleteOrder(ctx context.Context, req *oss.DeleteOrderRequest) (*oss.DeleteOrderResponse, error) {
	err := dao.DeleteOrder(req.OrderId)
	if err != nil {
		return &oss.DeleteOrderResponse{}, err
	}
	return &oss.DeleteOrderResponse{}, nil

}

func main() {
	dbConn.Connect()
	defer dbConn.Close()
	grpcserver.StartServer(func(server *grpc.Server) {
		oss.RegisterOrderStorageServiceServer(server, &service{})
	})
}
