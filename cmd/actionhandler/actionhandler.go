package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/handlers/createOrder", CreateOrder).Methods("POST")
	router.HandleFunc("/handlers/capturePayment", CapturePayment).Methods("POST")
	router.HandleFunc("/handlers/splitOrder", SplitOrder).Methods("POST")
	router.HandleFunc("/handlers/error", Error).Methods("POST")
	router.HandleFunc("/handlers/done", Done).Methods("POST")
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Print("CreateOrder")
}

func CapturePayment(w http.ResponseWriter, r *http.Request) {
	log.Print("CapturePayment")
}

func SplitOrder(w http.ResponseWriter, r *http.Request) {
	log.Print("SplitOrder")
}

func Error(w http.ResponseWriter, r *http.Request) {
	log.Print("Error")
}

func Done(w http.ResponseWriter, r *http.Request) {
	log.Print("Done")
}
