package main

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	"github.com/golang/protobuf/jsonpb"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestCRUD(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		t.Error("GOPATH undefined")
	}
	orderFile := goPath + "/src/github.com/albertsen/lessworkflow/data/test/order.json"
	data, err := ioutil.ReadFile(orderFile)
	if err != nil {
		t.Error(err)
	}
	newOrderJSON := string(data)
	var newOrder od.Order
	err = jsonpb.UnmarshalString(string(newOrderJSON), &newOrder)
	if err != nil {
		t.Error(err)
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		t.Error(err)
	}
	newOrder.Id = uuid.String()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	client := oss.NewOrderStorageServiceClient(conn)
	ctx := context.Background()
	getOrderResponse, err := client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: newOrder.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, getOrderResponse.Order)
	_, err = client.SaveOrder(ctx, &oss.SaveOrderRequest{Order: &newOrder})
	if err != nil {
		t.Error(err)
	}
	getOrderResponse, err = client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: newOrder.Id})
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, getOrderResponse.Order)
	newOrderJSON, err = new(jsonpb.Marshaler).MarshalToString(&newOrder)
	if err != nil {
		t.Error(err)
	}
	orderJSON, err := new(jsonpb.Marshaler).MarshalToString(getOrderResponse.Order)
	if err != nil {
		t.Error(err)
	}
	assert.JSONEq(t, newOrderJSON, orderJSON)
	_, err = client.DeleteOrder(ctx, &oss.DeleteOrderRequest{OrderId: newOrder.Id})
	if err != nil {
		t.Error(err)
	}
	getOrderResponse, err = client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: newOrder.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, getOrderResponse.Order)
}
