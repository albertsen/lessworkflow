package main

import (
	"context"
	"log"
	"os"
	"testing"

	pd "github.com/albertsen/lessworkflow/gen/proto/processdef"
	pds "github.com/albertsen/lessworkflow/gen/proto/processdefservice"
	"github.com/albertsen/lessworkflow/pkg/testing/cmpopts"
	tu "github.com/albertsen/lessworkflow/pkg/testing/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

var client pds.ProcessDefServiceClient
var ctx context.Context

func loadProcessDef(t *testing.T) *pd.ProcessDef {
	var processDef pd.ProcessDef
	err := tu.LoadTestData("process", &processDef)
	if err != nil {
		t.Error(err)
	}
	return &processDef
}

func TestMain(m *testing.M) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to Process Def Service: %s", err)
	}
	defer conn.Close()
	client = pds.NewProcessDefServiceClient(conn)
	ctx = context.Background()
	os.Exit(m.Run())
}

func TestCRUD(t *testing.T) {
	refProcessDef := loadProcessDef(t)
	createProcessDefResponse, err := client.CreateProcessDef(ctx, &pds.CreateProcessDefRequest{ProcessDef: refProcessDef})
	if err != nil {
		t.Error(err)
	}
	createdProcessDef := createProcessDefResponse.ProcessDef
	assert.Assert(t, createdProcessDef != nil, "CreateProcessDef didn't return process def")
	getProcessDefResponse, err := client.GetProcessDef(ctx, &pds.GetProcessDefRequest{ProcessDefId: refProcessDef.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, getProcessDefResponse.ProcessDef != nil, "GetProcessDef did not return processdef")
	assert.DeepEqual(t, refProcessDef, getProcessDefResponse.ProcessDef, cmpopts.IgnoreInternalProtbufFieldsOption)
	refProcessDef.Description = "My updated description"
	refProcessDef.Workflow.Actions["newAction"] = &pd.ActionDef{
		Handler: "handler",
		Transitions: map[string]string{
			"ok": "ok",
		},
	}
	_, err = client.UpdateProcessDef(ctx, &pds.UpdateProcessDefRequest{ProcessDef: refProcessDef})
	if err != nil {
		t.Error(err)
	}
	getProcessDefResponse, err = client.GetProcessDef(ctx, &pds.GetProcessDefRequest{ProcessDefId: refProcessDef.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, getProcessDefResponse.ProcessDef != nil, "GetProcessDef did not return process def")
	assert.DeepEqual(t, refProcessDef, getProcessDefResponse.ProcessDef, cmpopts.IgnoreInternalProtbufFieldsOption)
	_, err = client.DeleteProcessDef(ctx, &pds.DeleteProcessDefRequest{ProcessDefId: refProcessDef.Id})
	if err != nil {
		t.Error(err)
	}
	getProcessDefResponse, err = client.GetProcessDef(ctx, &pds.GetProcessDefRequest{ProcessDefId: refProcessDef.Id})
	if err != nil {
		t.Error(err)
	}
	assert.Assert(t, is.Nil(getProcessDefResponse.ProcessDef))
}

func TestCannotCreateProcessDefWithoutID(t *testing.T) {
	processDef := loadProcessDef(t)
	processDef.Id = ""
	_, err := client.CreateProcessDef(ctx, &pds.CreateProcessDefRequest{ProcessDef: processDef})
	assert.Equal(t, codes.InvalidArgument, tu.ErrToGRPCStatusCode(t, err))
}

func TestCannotGetOrderWithoutID(t *testing.T) {
	_, err := client.GetProcessDef(ctx, &pds.GetProcessDefRequest{})
	assert.Equal(t, codes.InvalidArgument, tu.ErrToGRPCStatusCode(t, err))
}

func TestCannotUpdateOrderWithoutID(t *testing.T) {
	processDef := loadProcessDef(t)
	processDef.Id = ""
	_, err := client.UpdateProcessDef(ctx, &pds.UpdateProcessDefRequest{ProcessDef: processDef})
	assert.Equal(t, codes.InvalidArgument, tu.ErrToGRPCStatusCode(t, err))
}

func TestCannotUpdateNonExistingOrder(t *testing.T) {
	processDef := loadProcessDef(t)
	processDef.Id = "doesNotExist"
	_, err := client.UpdateProcessDef(ctx, &pds.UpdateProcessDefRequest{ProcessDef: processDef})
	assert.Equal(t, codes.NotFound, tu.ErrToGRPCStatusCode(t, err))
}

func TestCannotDeleteOrderWithoutID(t *testing.T) {
	_, err := client.DeleteProcessDef(ctx, &pds.DeleteProcessDefRequest{})
	assert.Equal(t, codes.InvalidArgument, tu.ErrToGRPCStatusCode(t, err))
}
