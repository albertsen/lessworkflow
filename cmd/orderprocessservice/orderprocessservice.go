package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	httpUtil "github.com/albertsen/lessworkflow/pkg/net/http/util"
	"github.com/gorilla/mux"
)

var documentserviceURL string

type process struct {
	ProcessDefID string `json:"processDefId"`
	OrderID      string `json:"orderId"`
	ProcessID    string `json:"processId"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/processes/", StartProcess).Methods("POST")
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8000"
	}
	documentserviceURL = os.Getenv("DOCUMENTSERVICE_URL")
	if documentserviceURL == "" {
		documentserviceURL = "http://documentservice"
	}
	log.Printf("orderprocessservice listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func StartProcess(w http.ResponseWriter, r *http.Request) {
	var proc process
	if err := json.NewDecoder(r.Body).Decode(&proc); err != nil {
		httpUtil.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}
}
