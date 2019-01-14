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
	Document   *doc.Document
	Error      error
	StatusCode int
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

func getDocument(URL string, C chan result) {
	if URL == "" {
		C <- result{Error: errors.New("No URL provided for document to load")}
		return
	}
	res, err := httpClient.Get(URL)
	if err != nil {
		var statusCode int
		if res != nil {
			statusCode = res.StatusCode
		}
		C <- result{Error: err, StatusCode: statusCode}
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
		C <- result{Error: fmt.Errorf("Error parsing document from [%s]: %s", URL, err)}
	}
	C <- result{Document: &doc, Error: err}
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
	go getDocument(startProcReq.ProcessDefURL, processDefChan)
	go getDocument(startProcReq.DocumentURL, docChan)
	erroneousResults := make([]result, 0)
	for i := 0; i < 2; i++ {
		select {
		case processDefRes := <-processDefChan:
			if processDefRes.Error != nil {
				erroneousResults = append(erroneousResults, processDefRes)
			} else {
				action.ProcessDef = processDefRes.Document
				action.Name = "startProcess"
			}
		case docRes := <-docChan:
			if docRes.Error != nil {
				erroneousResults = append(erroneousResults, docRes)
			} else {
				action.Document = docRes.Document
			}
		}
	}
	if l := len(erroneousResults); l > 0 {
		errorMessages := make([]string, l)
		for i := 0; i < l; i++ {
			errorMessages[i] = erroneousResults[i].Error.Error()
		}
		httpUtil.SendError(w, http.StatusInternalServerError, strings.Join(errorMessages[:], ", "))
		return
	}
	httpUtil.SendOK(w, action)
}
