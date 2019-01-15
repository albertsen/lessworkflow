package main

import (
	"net/http"
	"testing"

	"github.com/albertsen/lessworkflow/pkg/data/process"
	"github.com/albertsen/lessworkflow/pkg/rest/client"

	"gotest.tools/assert"
)

var (
	processServiceURL = "http://localhost:8001/processes"
)

func TestCreateProcess(t *testing.T) {
	documentURL := "http://documentservice/documents/orders/neworder"
	processDefURL := "http://documentservice/documents/processdefs/fulfilorder"
	refProc := process.Process{
		DocumentURL:   documentURL,
		ProcessDefURL: processDefURL,
	}
	var createdProc process.Process
	res, err := client.Post(processServiceURL, &refProc, &createdProc)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf(string(res.Body))
	}
	assert.Assert(t, createdProc.ID != "", "In created process, ID should not be empty")
	assert.Assert(t, createdProc.TimeCreated != nil, "In created process, TimeCreated should not be nil")
	assert.Assert(t, createdProc.TimeUpdated != nil, "In created process, TimeUpdated should not be nil")
	assert.Equal(t, createdProc.Status, "CREATED")
	assert.Equal(t, createdProc.DocumentURL, documentURL)
	assert.Equal(t, createdProc.ProcessDefURL, processDefURL)
}
