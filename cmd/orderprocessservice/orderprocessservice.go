package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/processes/{type}", CreateDocument).Methods("POST")
	router.HandleFunc("/documents/{type}/{id}", GetDocument).Methods("GET")
	router.HandleFunc("/documents/{type}/{id}", UpdateDocument).Methods("PUT")
	router.HandleFunc("/documents/{type}/{id}", DeleteDocument).Methods("DELETE")
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	log.Printf("documentservice listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
