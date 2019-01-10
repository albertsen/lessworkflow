package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	httpUtil "github.com/albertsen/lessworkflow/pkg/net/http/util"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents/{type}", CreateDocument).Methods("POST")
	router.HandleFunc("/documents/{type}/{id}", GetDocument).Methods("GET")
	router.HandleFunc("/documents/{type}/{id}", UpdateDocument).Methods("PUT")
	router.HandleFunc("/documents/{type}/{id}", DeleteDocument).Methods("DELETE")
	dbConn.Connect()
	defer dbConn.Close()
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
		httpUtil.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if doc.Type == "" {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Document does not have a type. Should be: %s", docType))
		return
	}
	if doc.Type != docType {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Type of provided document [%s] does not match type provided in URL [%s]", doc.Type, docType))
		return
	}
	if doc.Version > 0 {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("New document cannot have a version > 0: %s", doc.Version))
	}
	if doc.ID == "" {
		if uuid, err := uuid.NewV4(); err != nil {
			httpUtil.SendError(w, http.StatusInternalServerError, err)
			return
		} else {
			doc.ID = uuid.String()
		}
	}
	now := time.Now()
	doc.TimeCreated = &now
	doc.TimeUpdated = &now
	doc.Version = 1
	if err := dbConn.DB().Insert(&doc); err != nil {
		pgError, ok := err.(pg.Error)
		if ok && pgError.IntegrityViolation() {
			httpUtil.SendError(w, http.StatusConflict, fmt.Sprintf("Document already exists with ID: %s", doc.ID))
		} else {
			httpUtil.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}
	httpUtil.SendResponse(w, http.StatusCreated, doc)
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	docType := mux.Vars(r)["type"]
	docID := mux.Vars(r)["id"]
	var doc doc.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if doc.Type == "" {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Document does not have a type. Should be: %s", docType))
		return
	}
	if doc.Type != docType {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Type of provided document [%s] does not match type provided in URL [%s]", doc.Type, docType))
		return
	}
	if doc.ID == "" {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Document does not have an ID. Should be: %s", docID))
		return
	}
	if doc.ID != docID {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("ID of provided document [%s] does not match ID provided in URL [%s]", doc.ID, docID))
		return
	}
	if doc.Version <= 0 {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, fmt.Sprintf("No or invalid version of document [%d]", doc.Version))
		return
	}
	now := time.Now()
	doc.TimeUpdated = &now
	if err := dbConn.DB().Update(&doc); err != nil {
		if err == pg.ErrNoRows {
			httpUtil.SendError(w, http.StatusNotFound, fmt.Sprintf("No document found with ID: %s", doc.ID))
		} else {
			httpUtil.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}
	httpUtil.SendOK(w, doc)
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	docID := mux.Vars(r)["id"]
	docType := mux.Vars(r)["type"]
	doc := &doc.Document{
		ID:   docID,
		Type: docType,
	}
	if err := dbConn.DB().Select(doc); err != nil {
		if err == pg.ErrNoRows {
			httpUtil.SendError(w, http.StatusNotFound, fmt.Sprintf("No document found with ID: %s", doc.ID))
		} else {
			httpUtil.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}
	httpUtil.SendOK(w, doc)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	docID := mux.Vars(r)["id"]
	docType := mux.Vars(r)["type"]
	doc := &doc.Document{
		ID:   docID,
		Type: docType,
	}
	if err := dbConn.DB().Delete(doc); err != nil {
		if err == pg.ErrNoRows {
			httpUtil.SendError(w, http.StatusNotFound, fmt.Sprintf("No document found with ID: %s", doc.ID))
		} else {
			httpUtil.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}
	httpUtil.SendOK(w, nil)
}
