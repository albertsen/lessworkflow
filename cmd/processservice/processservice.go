package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/albertsen/lessworkflow/pkg/data/action"
	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	pd "github.com/albertsen/lessworkflow/pkg/data/processdef"
	uuid "github.com/satori/go.uuid"

	httpUtil "github.com/albertsen/lessworkflow/pkg/net/http/util"
	"github.com/gorilla/mux"
)

type startProcessRequest struct {
	ProcessDefID string `json:"processDefId"`
	DocumentID   string `json:"documentId"`
	DocumentType string `json:"documentType"`
}

type startProceessResponse struct {
	ProcessID string `json:"processId`
}

type result struct {
	Value interface{}
	Error error
	Type  string
}

var (
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	documentserviceURL string
	typeProcessDef     = "processdefs"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/processes", StartProcess).Methods("POST")
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	documentserviceURL = os.Getenv("DOCUMENTSERVICE_URL")
	if documentserviceURL == "" {
		documentserviceURL = "http://documentservice"
	}
	log.Printf("processservice listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func loadDocument(DocType string, DocID string) (*doc.Document, error) {
	res, err := httpClient.Get(documentserviceURL + "/documents/" + url.PathEscape(DocType) + "/" + url.PathEscape(DocID))
	if err != nil {
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil
	}
	var doc doc.Document
	err = json.Unmarshal(data, &doc)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing document of type [%s] with ID [%s]: %s", DocType, DocID, err))
	}
	return &doc, nil
}

func getDocument(DocType string, DocID string, c chan result) {
	doc, err := loadDocument(DocType, DocID)
	c <- result{Value: doc, Error: err, Type: DocType}
}

func getProcessDef(ProcessDefId string, c chan result) {
	doc, err := loadDocument(typeProcessDef, ProcessDefId)
	if err != nil {
		c <- result{Error: err}
		return
	}
	if doc == nil {
		c <- result{}
		return
	}
	var processDef pd.ProcessDef
	err = json.Unmarshal(doc.Content, &processDef)
	log.Println(processDef)
	if err != nil {
		c <- result{Error: errors.New(fmt.Sprintf("Error parsing process definition: %s", err))}
		return
	}
	c <- result{Value: &processDef, Type: typeProcessDef}
}

func StartProcess(w http.ResponseWriter, r *http.Request) {
	var startProcReq startProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&startProcReq); err != nil {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if startProcReq.DocumentID == "" || startProcReq.DocumentType == "" || startProcReq.ProcessDefID == "" {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, "Invalid request: Needs all of documentType, documentId, processDefId")
		return
	}
	c := make(chan result, 2)
	// Make two parallel calls to the document service
	go getDocument(startProcReq.DocumentType, startProcReq.DocumentID, c)
	go getProcessDef(startProcReq.ProcessDefID, c)
	var results [2]result
	results[0], results[1] = <-c, <-c
	errors := make([]error, 0)
	var processDef *pd.ProcessDef
	var document *doc.Document
	for _, r := range results {
		if r.Error == nil {
			if r.Type == typeProcessDef {
				processDef = r.Value.(*pd.ProcessDef)
			} else {
				document = r.Value.(*doc.Document)
			}
		} else {
			errors = append(errors, r.Error)
		}
	}
	if len(errors) > 0 {
		httpUtil.SendError(w, http.StatusInternalServerError, errors)
		return
	}
	if document == nil || processDef == nil {
		var errorMessage string
		if document == nil {
			errorMessage += fmt.Sprintf("No document found of type [%s] with ID [%s]. ", startProcReq.DocumentType, startProcReq.DocumentID)
		}
		if processDef == nil {
			errorMessage += fmt.Sprintf("No process definotion found with ID: %s. ", startProcReq.ProcessDefID)
		}
		httpUtil.SendError(w, http.StatusNotFound, errorMessage)
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
		ProcessID:  processID,
		Name:       processDef.Workflow.Start,
		ProcessDef: processDef,
		Document:   document,
	}
	httpUtil.SendOK(w, action)

}
