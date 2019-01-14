package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/albertsen/lessworkflow/pkg/data/action"
	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	pd "github.com/albertsen/lessworkflow/pkg/data/processdef"
	uuid "github.com/satori/go.uuid"

	httpUtil "github.com/albertsen/lessworkflow/pkg/net/http/util"
	"github.com/gorilla/mux"
)

type startProcessRequest struct {
	DocumentURL   string `json:"documentURL"`
	ProcessDefURL string `json:"processDefURL"`
}

type startProceessResponse struct {
	ProcessURL string `json:"processId"`
}

type result struct {
	Value interface{}
	Error error
}

var (
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/processes", StartProcess).Methods("POST")
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	log.Printf("processservice listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func getDocument(URL string, C chan result, Converter func(*doc.Document) (interface{}, error)) {
	if URL == "" {
		C <- result{Error: errors.New("No URL provided for document to load")}
		return
	}
	res, err := httpClient.Get(URL)
	if err != nil {
		if res != nil && res.StatusCode == http.StatusNotFound {
			C <- result{}
			return
		}
		C <- result{Error: err}
		return
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		C <- result{}
		return
	}
	var doc doc.Document
	err = json.Unmarshal(data, &doc)
	if err != nil {
		C <- result{Error: errors.New(fmt.Sprintf("Error parsing document from [%s]: %s", URL, err))}
	}
	value, err := Converter(&doc)
	C <- result{Value: value, Error: err}
}

func StartProcess(w http.ResponseWriter, r *http.Request) {
	var startProcReq startProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&startProcReq); err != nil {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if startProcReq.ProcessDefURL == "" || startProcReq.DocumentURL == "" {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, "Invalid request: Needs all of processDefURL and documentURL")
		return
	}
	var processID string
	if uuid, err := uuid.NewV4(); err != nil {
		httpUtil.SendError(w, http.StatusInternalServerError, err)
		return
	} else {
		processID = uuid.String()
	}
	action := action.Action{
		ProcessID: processID,
	}
	processDefChan := make(chan result)
	docChan := make(chan result)
	// Make two parallel calls to the document service
	go getDocument(startProcReq.ProcessDefURL, processDefChan, func(Doc *doc.Document) (interface{}, error) {
		return Doc, nil
	})
	go getDocument(startProcReq.ProcessDefURL, docChan, func(Doc *doc.Document) (interface{}, error) {
		var processDef pd.ProcessDef
		if err := json.Unmarshal(Doc.Content, &processDef); err != nil {
			return nil, err
		}
		return processDef, nil
	})
	errorMessages := make([]string, 0)
	notFound := make([]string, 0)
	for i := 0; i < 2; i++ {
		select {
		case processDefRes := <-processDefChan:
			if processDefRes.Error != nil {
				errorMessages = append(errorMessages, processDefRes.Error.Error())
			} else {
				processDef := processDefRes.Value.(*pd.ProcessDef)
				if processDef == nil {
					notFound = append(errorMessages, fmt.Sprintf("Process definition not found at: %s", startProcReq.ProcessDefURL))
				} else {
					action.ProcessDef = processDef
					action.Name = processDef.Workflow.Start
				}
			}
		case docRes := <-docChan:
			if docRes.Error != nil {
				errorMessages = append(errorMessages, docRes.Error.Error())
			} else {
				if action.Document == nil {
					notFound = append(errorMessages, fmt.Sprintf("Document not found at: %s", startProcReq.DocumentURL))
				}
				action.Document = docRes.Value.(*doc.Document)
			}
		}
	}
	if len(errorMessages) > 0 {
		httpUtil.SendError(w, http.StatusInternalServerError, strings.Join(errorMessages[:], ", "))
		return
	}
	if len(notFound) > 0 {
		httpUtil.SendError(w, http.StatusNotFound, strings.Join(notFound[:], ", "))
		return
	}
	httpUtil.SendOK(w, action)
}
