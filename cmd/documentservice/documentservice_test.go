package main

import (
	"net/http"
	"testing"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	od "github.com/albertsen/lessworkflow/pkg/data/order"
	"github.com/albertsen/lessworkflow/pkg/rest/client"
	"github.com/google/go-cmp/cmp/cmpopts"

	tu "github.com/albertsen/lessworkflow/pkg/testing/utils"
	"gotest.tools/assert"
)

var (
	documentServiceURL = "http://localhost:8000/documents"
)

func TestCRUD(t *testing.T) {
	var refDoc doc.Document
	err := tu.LoadData("sample/order.json", &refDoc)
	if err != nil {
		t.Error(err)
	}
	var createdDoc doc.Document
	res, err := client.Post(documentServiceURL+"/"+refDoc.Type, &refDoc, &createdDoc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusCreated, "HTTP status should be CREATED")
	assert.Assert(t, createdDoc.ID != "", "In created document, ID should not be empty")
	assert.Assert(t, createdDoc.TimeCreated != nil, "In created document, TimeCreated should not be nil")
	assert.Assert(t, createdDoc.TimeUpdated != nil, "In created document, TimeUpdated should not be nil")
	assert.Assert(t, createdDoc.Version <= 1, "In created document, Version should not be greater than 1")
	refDoc.ID = createdDoc.ID
	refDoc.TimeCreated = createdDoc.TimeCreated
	refDoc.TimeUpdated = createdDoc.TimeUpdated
	refDoc.Version = createdDoc.Version
	docURL := documentServiceURL + "/" + refDoc.Type + "/" + refDoc.ID
	var storedDoc doc.Document
	res, err = client.Get(docURL, &storedDoc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, "HTTP status should be OK")
	assertDocsEqual(t, &refDoc, &storedDoc)
	var order od.Order
	err = refDoc.GetContent(&order)
	if err != nil {
		t.Error(err)
	}
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
	order.LineItems = append(order.LineItems, lineItem)
	order.TotalPrice.Value += lineItem.TotalPrice.Value
	err = refDoc.SetContent(&order)
	if err != nil {
		t.Error(err)
	}
	res, err = client.Put(docURL, refDoc, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, "HTTP status should be OK")
	res, err = client.Get(docURL, &storedDoc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, "HTTP status should be OK")
	assert.Assert(t, storedDoc.TimeUpdated.After(*refDoc.TimeUpdated), "TimeUpdated not updated")
	refDoc.TimeUpdated = storedDoc.TimeUpdated
	assertDocsEqual(t, &refDoc, &storedDoc)
	res, err = client.Delete(docURL)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK, "HTTP status should be OK")
	res, err = client.Get(docURL, &storedDoc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusNotFound, "HTTP status should be NOT FOUND")
}

func assertDocsEqual(T *testing.T, Doc1 *doc.Document, Doc2 *doc.Document) {
	assert.Equal(T, Doc1.ID, Doc2.ID)
	assert.Assert(T, Doc1.TimeCreated.Equal(*Doc2.TimeCreated))
	assert.Assert(T, Doc1.TimeUpdated.Equal(*Doc2.TimeUpdated))
	assert.Equal(T, Doc1.Type, Doc2.Type)
	assert.Equal(T, Doc1.Version, Doc2.Version)
	var order1 od.Order
	var order2 od.Order
	err := Doc1.GetContent(&order1)
	if err != nil {
		T.Error(err)
	}
	err = Doc2.GetContent(&order2)
	if err != nil {
		T.Error(err)
	}
	assert.DeepEqual(T, order1, order2, cmpopts.IgnoreUnexported(doc.Document{}))
}
