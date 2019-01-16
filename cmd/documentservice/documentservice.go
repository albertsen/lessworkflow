package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	"github.com/albertsen/lessworkflow/pkg/db/repo"
	"github.com/albertsen/lessworkflow/pkg/rest/server"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents/{type}", CreateDocument).Methods("POST")
	router.HandleFunc("/documents/{type}/{id}", GetDocument).Methods("GET")
	router.HandleFunc("/documents/{type}/{id}", UpdateDocument).Methods("PUT")
	router.HandleFunc("/documents/{type}/{id}", DeleteDocument).Methods("DELETE")
	repo.Connect()
	defer repo.Close()
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	log.Printf("documentservice listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	docType := mux.Vars(r)["type"]
	var doc doc.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		server.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if doc.Type == "" {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Document does not have a type. Should be: %s", docType))
		return
	}
	if doc.Type != docType {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Type of provided document [%s] does not match type provided in URL [%s]", doc.Type, docType))
		return
	}
	if doc.Version > 0 {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("New document cannot have a version > 0: %d", doc.Version))
	}
	if doc.ID == "" {
		if uuid, err := uuid.NewV4(); err != nil {
			server.SendError(w, http.StatusInternalServerError, err)
			return
		} else {
			doc.ID = uuid.String()
		}
	}
	now := time.Now().Truncate(time.Microsecond)
	doc.TimeCreated = &now
	doc.TimeUpdated = &now
	doc.Version = 1
	statusCode, err := repo.Insert(&doc)
	server.SendResponseOrError(w, http.StatusCreated, statusCode, doc, err)
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	docType := mux.Vars(r)["type"]
	docID := mux.Vars(r)["id"]
	var doc doc.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		server.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if doc.Type == "" {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Document does not have a type. Should be: %s", docType))
		return
	}
	if doc.Type != docType {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Type of provided document [%s] does not match type provided in URL [%s]", doc.Type, docType))
		return
	}
	if doc.ID == "" {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Document does not have an ID. Should be: %s", docID))
		return
	}
	if doc.ID != docID {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("ID of provided document [%s] does not match ID provided in URL [%s]", doc.ID, docID))
		return
	}
	if doc.Version <= 0 {
		server.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("No or invalid version of document [%d]", doc.Version))
		return
	}
	now := time.Now().Truncate(time.Microsecond)
	doc.TimeUpdated = &now
	statusCode, err := repo.Update(&doc)
	server.SendResponseOrError(w, http.StatusOK, statusCode, doc, err)
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	docID := mux.Vars(r)["id"]
	docType := mux.Vars(r)["type"]
	doc := doc.Document{
		ID:   docID,
		Type: docType,
	}
	statusCode, err := repo.Select(&doc)
	server.SendResponseOrError(w, http.StatusOK, statusCode, doc, err)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	docID := mux.Vars(r)["id"]
	docType := mux.Vars(r)["type"]
	doc := doc.Document{
		ID:   docID,
		Type: docType,
	}
	statusCode, err := repo.Delete(&doc)
	server.SendResponseOrError(w, http.StatusOK, statusCode, doc, err)
}
