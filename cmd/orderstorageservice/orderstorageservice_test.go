package main

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	pb "github.com/albertsen/lessworkflow/gen/proto/order"
	"github.com/golang/protobuf/jsonpb"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestCRUD(t *testing.T) {
	goPath, _ := os.LookupEnv("GOPATH")
	orderFile := goPath + "/src/github.com/albertsen/lessworkflow/data/test/order.json"
	data, err := ioutil.ReadFile(orderFile)
	if err != nil {
		t.Error(err)
	}
	newOrderJSON := string(data)
	var newOrder pb.Order
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
	client := pb.NewOrderStorageServiceClient(conn)
	ctx := context.Background()
	_, err = client.SaveOrder(ctx, &newOrder)
	if err != nil {
		t.Error(err)
	}
	order, err := client.GetOrder(ctx, &pb.GetOrderRequest{OrderId: newOrder.Id})
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, order)
	newOrderJSON, err = new(jsonpb.Marshaler).MarshalToString(&newOrder)
	if err != nil {
		t.Error(err)
	}
	orderJSON, err := new(jsonpb.Marshaler).MarshalToString(order)
	if err != nil {
		t.Error(err)
	}
	assert.JSONEq(t, newOrderJSON, orderJSON)
}
