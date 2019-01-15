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
	Document *doc.Document
	Err      error
	Response *client.Response
}

func (r *result) IsError() bool {
	return r.Err != nil || r.Response.StatusCode != http.StatusOK
}

func (r *result) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	if r.Response.Message != "" {
		return r.Response.Message
	}
	return string(r.Response.Body)
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
	res, err := client.Get(URL, &doc)
	if err != nil {
		C <- result{Err: err, Response: res}
	} else {
		C <- result{Document: &doc, Response: res}
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
	errorMessages := make([]string, 0)
	for i := 0; i < 2; i++ {
		select {
		case processDefRes := <-processDefChan:
			if !processDefRes.IsError() {
				action.ProcessDef = processDefRes.Document
				action.Name = "startProcess"
			} else {
				errorMessages = append(errorMessages, processDefRes.Error())
			}
		case docRes := <-docChan:
			if !docRes.IsError() {
				action.Document = docRes.Document
			} else {
				errorMessages = append(errorMessages, docRes.Error())
			}
		}
	}
	if l := len(errorMessages); l > 0 {
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
