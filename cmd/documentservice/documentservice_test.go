package main

import (
	"context"
	"log"
	"os"
	"testing"

	doc "github.com/albertsen/lessworkflow/gen/proto/document"
	ds "github.com/albertsen/lessworkflow/gen/proto/documentservice"
	"github.com/albertsen/lessworkflow/gen/proto/order"
	"github.com/albertsen/lessworkflow/pkg/testing/cmpopts"
	tu "github.com/albertsen/lessworkflow/pkg/testing/utils"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

var client ds.DocumentServiceClient
var ctx context.Context

func loadOrder(t *testing.T) *order.Order {
	var order order.Order
	err := tu.LoadTestData("order", &order)
	if err != nil {
		t.Error(err)
	}
	return &order
}

func TestMain(m *testing.M) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to Order Storage Service: %s", err)
	}
	defer conn.Close()
	client = ds.NewDocumentServiceClient(conn)
	ctx = context.Background()
	os.Exit(m.Run())
}

func TestCRUD(t *testing.T) {
	refOrder := loadOrder(t)
	data, err := ptypes.MarshalAny(refOrder)
	if err != nil {
		t.Error(err)
	}
	refDoc := doc.Document{
		Id:   "newOrder",
		Type: "order",
		Data: data,
	}
	createDocResponse, err := client.CreateDocument(ctx, &ds.CreateDocumentRequest{Document: &refDoc})
	if err != nil {
		t.Error(err)
	}
	createdDoc := createDocResponse.Document
	assert.Assert(t, createdDoc != nil, "CreateDocument didn't return document")
	assert.Assert(t, createdDoc.Id != "", "In created document, Id should not be empty")
	assert.Assert(t, createdDoc.TimeCreated != nil, "In created document, TimeCreated should not be nil")
	assert.Assert(t, createdDoc.TimeUpdated != nil, "In created document, TimeUpdated should not be nil")
	assert.Assert(t, createdDoc.Version > 0, "In created document, Version should not be greater than 0")
	refDoc.Id = createdDoc.Id
	refDoc.TimeCreated = createdDoc.TimeCreated
	refDoc.TimeUpdated = createdDoc.TimeUpdated
	refDoc.Version = createdDoc.Version
	getDocResponse, err := client.GetDocument(ctx, &ds.GetDocumentRequest{DocumentId: createDocResponse.Document.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, getDocResponse.Document != nil, "GetDociment did not return order")
	assert.DeepEqual(t, refDoc, getDocResponse.Document, cmpopts.IgnoreInternalProtbufFieldsOption)
	lineItem := &order.LineItem{
		Count: 3,
		ItemPrice: &order.MonetaryAmount{
			Value:    100,
			Currency: "EUR",
		},
		ProductId:          "oettinger",
		ProductDescription: "Oettinger Bier",
		TotalPrice: &order.MonetaryAmount{
			Value:    300,
			Currency: "EUR",
		},
	}
	refOrder.LineItems = append(refOrder.LineItems, lineItem)
	refOrder.TotalPrice.Value += lineItem.TotalPrice.Value
	data, err = ptypes.MarshalAny(refOrder)
	if err != nil {
		t.Error(err)
	}
	refDoc.Data = data
	_, err = client.UpdateDocument(ctx, &ds.UpdateDocumentRequest{Document: &refDoc})
	if err != nil {
		t.Error(err)
	}
	getDocResponse, err = client.GetDocument(ctx, &ds.GetDocumentRequest{DocumentId: refDoc.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, getDocResponse.Document != nil, "GetDocument did not return order")
	assert.Assert(t, getDocResponse.Document.TimeUpdated.Nanos > refDoc.TimeUpdated.Nanos, "TimeUpdated not updated")
	refDoc.TimeUpdated = getDocResponse.Document.TimeUpdated
	assert.DeepEqual(t, refDoc, getDocResponse.Document, cmpopts.IgnoreInternalProtbufFieldsOption)
	_, err = client.DeleteDocument(ctx, &ds.DeleteDocumentRequest{DocumentId: refDoc.Id})
	if err != nil {
		t.Error(err)
	}
	getDocResponse, err = client.GetDocument(ctx, &ds.GetDocumentRequest{DocumentId: refDoc.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, is.Nil(getDocResponse.Document))
}

// func TestCannotCreateOrderWithID(t *testing.T) {
// 	order := loadOrder(t)
// 	order.Id = "anIDThatCannotBe"
// 	_, err := client.CreateOrder(ctx, &oss.CreateOrderRequest{Order: order})
// 	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
// }

// func TestCannotGetOrderWithoutID(t *testing.T) {
// 	_, err := client.GetOrder(ctx, &oss.GetOrderRequest{})
// 	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
// }

// func TestCannotUpdateOrderWithoutID(t *testing.T) {
// 	order := loadOrder(t)
// 	_, err := client.UpdateOrder(ctx, &oss.UpdateOrderRequest{Order: order})
// 	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
// }

// func TestCannotUpdateNonExistingOrder(t *testing.T) {
// 	order := loadOrder(t)
// 	uuid, err := uuid.NewV4()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	order.Id = uuid.String()
// 	_, err = client.UpdateOrder(ctx, &oss.UpdateOrderRequest{Order: order})
// 	assert.Equal(t, codes.NotFound, errToGRPCStatusCode(t, err))
// }

// func TestCannotDeleteOrderWithoutID(t *testing.T) {
// 	_, err := client.DeleteOrder(ctx, &oss.DeleteOrderRequest{})
// 	assert.Equal(t, codes.InvalidArgument, errToGRPCStatusCode(t, err))
// }

func errToGRPCStatusCode(t *testing.T, err error) codes.Code {
	stat, ok := status.FromError(err)
	assert.Assert(t, ok, "Error is not a gRPC error")
	return stat.Code()
}
