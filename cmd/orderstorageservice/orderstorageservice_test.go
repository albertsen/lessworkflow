package main

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func loadOrder() (*od.Order, error) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return nil, errors.New("GOPATH undefined")
	}
	orderFile := goPath + "/src/github.com/albertsen/lessworkflow/data/test/order.json"
	data, err := ioutil.ReadFile(orderFile)
	if err != nil {
		return nil, err
	}
	orderJSON := string(data)
	var order od.Order
	err = jsonpb.UnmarshalString(string(orderJSON), &order)
	if err != nil {
		return nil, err
	}
	return &order, err
}

func TestCRUD(t *testing.T) {
	newOrder, err := loadOrder()
	if err != nil {
		t.Error(err)
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	client := oss.NewOrderStorageServiceClient(conn)
	ctx := context.Background()
	createOrderResponse, err := client.CreateOrder(ctx, &oss.CreateOrderRequest{Order: newOrder})
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, createOrderResponse.OrderId)
	getOrderResponse, err := client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: createOrderResponse.OrderId})
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, getOrderResponse.Order)
	_, err = client.DeleteOrder(ctx, &oss.DeleteOrderRequest{OrderId: createOrderResponse.OrderId})
	if err != nil {
		t.Error(err)
	}
	getOrderResponse, err = client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: createOrderResponse.OrderId})
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, getOrderResponse.Order)
}
