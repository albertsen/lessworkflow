package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/albertsen/lessworkflow/pkg/data/action"
	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	"github.com/albertsen/lessworkflow/pkg/data/process"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	uuid "github.com/satori/go.uuid"

	"github.com/albertsen/lessworkflow/pkg/rest/client"
	"github.com/albertsen/lessworkflow/pkg/rest/server"
	"github.com/gorilla/mux"
)

type result struct {
	Document   *doc.Document
	Error      error
	StatusCode int
}

func main() {
	router := mux.NewRouter()
	dbConn.Connect()
	defer dbConn.Close()
	router.HandleFunc("/processes", StartProcess).Methods("POST")
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	log.Printf("processservice listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func getDocument(URL string, C chan result) {
	var doc doc.Document
	statusCode, err := client.Get(URL, &doc)
	if err != nil {
		C <- result{Error: err, StatusCode: statusCode}
	} else {
		C <- result{Document: &doc, StatusCode: statusCode}
	}
}

func StartProcess(w http.ResponseWriter, r *http.Request) {
	var process process.Process
	if err := json.NewDecoder(r.Body).Decode(&process); err != nil {
		server.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if process.ProcessDefURL == "" || process.DocumentURL == "" {
		server.SendError(w, http.StatusUnprocessableEntity, "Invalid request: Needs all of processDefURL and documentURL")
		return
	}
	var processID string
	if uuid, err := uuid.NewV4(); err != nil {
		server.SendError(w, http.StatusInternalServerError, err)
		return
	} else {
		processID = uuid.String()
	}
	process.ID = processID
	action := action.Action{
		ProcessID: processID,
	}
	processDefChan := make(chan result)
	docChan := make(chan result)
	// Make two parallel calls to the document service
	go getDocument(process.ProcessDefURL, processDefChan)
	go getDocument(process.DocumentURL, docChan)
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
		server.SendError(w, http.StatusInternalServerError, strings.Join(errorMessages[:], ", "))
		return
	}
	now := time.Now().Truncate(time.Microsecond)
	process.TimeCreated = &now
	process.TimeUpdated = &now
	process.Status = "CREATED"
	err := dbConn.DB().Insert(&process)
	if err != nil {
		server.SendError(w, http.StatusInternalServerError, err)
		return
	}
	server.SendResponse(w, http.StatusCreated, process)
}
