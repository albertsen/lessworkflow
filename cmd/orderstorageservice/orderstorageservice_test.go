package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	"github.com/albertsen/lessworkflow/pkg/testing/cmpopts"
	"github.com/golang/protobuf/jsonpb"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

var client oss.OrderStorageServiceClient
var ctx context.Context

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

func TestMain(m *testing.M) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to Order Storage Service: %s", err)
	}
	defer conn.Close()
	client = oss.NewOrderStorageServiceClient(conn)
	ctx = context.Background()
	os.Exit(m.Run())
}

func TestCRUD(t *testing.T) {
	refOrder, err := loadOrder()
	if err != nil {
		t.Error(err)
	}
	createOrderResponse, err := client.CreateOrder(ctx, &oss.CreateOrderRequest{Order: refOrder})
	if err != nil {
		t.Error(err)
	}
	createdOrder := createOrderResponse.Order
	assert.Assert(t, createdOrder != nil, "CreateOrder didn't return order")
	assert.Assert(t, createdOrder.Id != "", "In created order, Id should not be empty")
	assert.Assert(t, createdOrder.TimeCreated != nil, "In created order, TimeCreated should not be nil")
	assert.Assert(t, createdOrder.TimePlaced != nil, "In created order, TimePlaced should not be nil")
	assert.Assert(t, createdOrder.TimeUpdated != nil, "In created order, TimeUpdated should not be nil")
	assert.Assert(t, createdOrder.Status != "", "In created order, Status should not be empty")
	assert.Assert(t, createdOrder.Version > 0, "In created order, Version should not be greater than 0")
	refOrder.Id = createdOrder.Id
	refOrder.TimeCreated = createdOrder.TimeCreated
	refOrder.TimePlaced = createdOrder.TimePlaced
	refOrder.TimeUpdated = createdOrder.TimeUpdated
	refOrder.Status = createdOrder.Status
	refOrder.Version = createdOrder.Version
	getOrderResponse, err := client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: createOrderResponse.Order.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, getOrderResponse.Order != nil, "GetOrder did not return order")
	assert.DeepEqual(t, refOrder, getOrderResponse.Order, cmpopts.IgnoreInternalProtbufFieldsOption)
	lineItem := &od.LineItem{
		Count: 3,
		ItemPrice: &od.MonetaryAmount{
			Value:    100,
			Currency: "EUR",
		},
		ProductId:          "oettinger",
		ProductDescription: "Oettinger Bier",
		TotalPrice: &od.MonetaryAmount{
			Value:    300,
			Currency: "EUR",
		},
	}
	refOrder.Details.Ordered.LineItems = append(refOrder.Details.Ordered.LineItems, lineItem)
	refOrder.Details.Ordered.TotalPrice.Value += lineItem.TotalPrice.Value
	_, err = client.UpdateOrder(ctx, &oss.UpdateOrderRequest{Order: refOrder})
	if err != nil {
		t.Error(err)
	}
	getOrderResponse, err = client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: createOrderResponse.Order.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, getOrderResponse.Order != nil, "GetOrder did not return order")
	assert.Assert(t, getOrderResponse.Order.TimeUpdated.Nanos > refOrder.TimeUpdated.Nanos, "TimeUpdated not updated")
	refOrder.TimeUpdated = getOrderResponse.Order.TimeUpdated
	assert.DeepEqual(t, refOrder, getOrderResponse.Order, cmpopts.IgnoreInternalProtbufFieldsOption)
	_, err = client.DeleteOrder(ctx, &oss.DeleteOrderRequest{OrderId: createOrderResponse.Order.Id})
	if err != nil {
		t.Error(err)
	}
	getOrderResponse, err = client.GetOrder(ctx, &oss.GetOrderRequest{OrderId: createOrderResponse.Order.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, is.Nil(getOrderResponse.Order))
}

func TestCannotCreateOrderWithID(t *testing.T) {
	order, err := loadOrder()
	if err != nil {
		t.Error(err)
	}
	order.Id = "anIDThatCannotBe"
	_, err = client.CreateOrder(ctx, &oss.CreateOrderRequest{Order: order})
	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
}

func TestCannotGetOrderWithoutID(t *testing.T) {
	_, err := client.GetOrder(ctx, &oss.GetOrderRequest{})
	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
}

func TestCannotUpdateOrderWithoutID(t *testing.T) {
	order, err := loadOrder()
	if err != nil {
		t.Error(err)
	}
	_, err = client.UpdateOrder(ctx, &oss.UpdateOrderRequest{Order: order})
	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
}

func TestCannotUpdateNonExistingOrder(t *testing.T) {
	order, err := loadOrder()
	if err != nil {
		t.Error(err)
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		t.Error(err)
	}
	order.Id = uuid.String()
	_, err = client.UpdateOrder(ctx, &oss.UpdateOrderRequest{Order: order})
	assert.Equal(t, codes.NotFound, errToGRPCStatusCode(t, err))
}

func TestCannotDeleteOrderWithoutID(t *testing.T) {
	_, err := client.DeleteOrder(ctx, &oss.DeleteOrderRequest{})
	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
}

func errToGRPCStatusCode(t *testing.T, err error) codes.Code {
	stat, ok := status.FromError(err)
	assert.Assert(t, ok, "Error is not a gRPC error")
	return stat.Code()
}
