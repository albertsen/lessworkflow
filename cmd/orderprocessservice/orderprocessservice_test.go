package main

import (
	"context"
	"testing"

	ad "github.com/albertsen/lessworkflow/gen/proto/actiondata"
	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	oss "github.com/albertsen/lessworkflow/gen/proto/orderstorageservice"
	pd "github.com/albertsen/lessworkflow/gen/proto/processdef"
	tu "github.com/albertsen/lessworkflow/pkg/testing/utils"
	"github.com/gogo/protobuf/proto"
)

var client oss.OrderStorageServiceClient
var ctx context.Context

func loadOrder(t *testing.T) *od.Order {
	var order od.Order
	err := tu.LoadTestData("order", &order)
	if err != nil {
		t.Error(err)
	}
	return &order
}

func loadProcessDef(t *testing.T) *pd.ProcessDef {
	var processDef pd.ProcessDef
	err := tu.LoadTestData("process", &processDef)
	if err != nil {
		t.Error(err)
	}
	return &processDef
}

// func TestMain(m *testing.M) {
// 	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Cannot connect to Order Storage Service: %s", err)
// 	}
// 	defer conn.Close()
// 	client = oss.NewOrderStorageServiceClient(conn)
// 	ctx = context.Background()
// 	os.Exit(m.Run())
// }

func TestActionSize(t *testing.T) {
	order := loadOrder(t)
	processDef := loadProcessDef(t)
	action := ad.Action{
		Name:       "placeOrder",
		RetryCount: 1,
		Order:      order,
		ProcessDef: processDef,
	}
	data, err := proto.Marshal(&action)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Size of action %d", len(data))
}
