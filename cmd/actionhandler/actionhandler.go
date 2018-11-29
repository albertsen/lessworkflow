package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/albertsen/lessworkflow/pkg/action"
	"github.com/albertsen/lessworkflow/pkg/order"
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
	updateOrderStatus(w, r, "CREATED")
}

func CapturePayment(w http.ResponseWriter, r *http.Request) {
	updateOrderStatus(w, r, "PAYMENT_CAPTURED")
}

func SplitOrder(w http.ResponseWriter, r *http.Request) {
	updateOrderStatus(w, r, "ORDER_SPLIT")
}

func Error(w http.ResponseWriter, r *http.Request) {
	updateOrderStatus(w, r, "ERROR")

}

func Done(w http.ResponseWriter, r *http.Request) {
	updateOrderStatus(w, r, "DONE")
}

func updateOrderStatus(w http.ResponseWriter, r *http.Request, Result string) {
	log.Printf("Returning result: %s", Result)
	var order order.Order
	json.NewDecoder(r.Body).Decode(&order)
	order.Status = Result
	var actionResponse action.Response
	actionResponse.Result = Result
	actionResponse.Payload = action.Payload{
		ID:      order.ID,
		Type:    "order",
		Content: order,
	}
	json.NewEncoder(w).Encode(actionResponse)
}
