package main

import (
	"fmt"
	"net/http"
	"testing"

	od "github.com/albertsen/lessworkflow/pkg/data/order"
	"github.com/albertsen/lessworkflow/pkg/rest/client"

	tu "github.com/albertsen/lessworkflow/pkg/testing/utils"
	"gotest.tools/assert"
)

var (
	documentServiceURL = "http://localhost:8000"
)

func TestCRUD(t *testing.T) {
	var refOrder od.Order
	err := tu.LoadData("sample/order.json", &refOrder)
	if err != nil {
		t.Error(err)
	}
	var createdOrder od.Order
	res, err := client.Post(documentServiceURL+"/orders", &refOrder, &createdOrder)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusCreated, res.Message)
	assert.Assert(t, createdOrder.ID != "", "In created document, ID should not be empty")
	assert.Assert(t, createdOrder.TimeCreated != nil, "In created document, TimeCreated should not be nil")
	assert.Assert(t, createdOrder.TimeUpdated != nil, "In created document, TimeUpdated should not be nil")
	refOrder.ID = createdOrder.ID
	refOrder.TimeCreated = createdOrder.TimeCreated
	refOrder.TimeUpdated = createdOrder.TimeUpdated
	orderURL := documentServiceURL + "/orders/" + refOrder.ID
	var storedOrder od.Order
	res, err = client.Get(orderURL, &storedOrder)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, fmt.Sprintf("HTTP status should be OK - %s", res))
	assert.DeepEqual(t, &refOrder, &storedOrder)
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
	refOrder.LineItems = append(refOrder.LineItems, lineItem)
	refOrder.TotalPrice.Value += lineItem.TotalPrice.Value
	res, err = client.Put(orderURL, refOrder, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, fmt.Sprintf("HTTP status should be OK - %s", res))
	res, err = client.Get(orderURL, &storedOrder)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, fmt.Sprintf("HTTP status should be OK - %s", res))
	assert.Assert(t, storedOrder.TimeUpdated.After(*refOrder.TimeUpdated), "TimeUpdated not updated")
	refOrder.TimeUpdated = storedOrder.TimeUpdated
	assert.DeepEqual(t, &refOrder, &storedOrder)
	res, err = client.Delete(orderURL)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, fmt.Sprintf("HTTP status should be OK - %s", res))
	res, err = client.Get(orderURL, &storedOrder)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusNotFound, fmt.Sprintf("HTTP status should be NOT FOUND - %s", res))
}
