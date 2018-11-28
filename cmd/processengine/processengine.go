package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/albertsen/lessworkflow/pkg/msg"
	"github.com/albertsen/lessworkflow/pkg/process"
)

var (
	url         = flag.String("s", "nats://localhost:4222", "URL of messaging server.")
	topic       = flag.String("t", "actions", "Message topic.")
	processFile = flag.String("p", "", "Process descriptor file. Mandatory.")
	help        = flag.Bool("h", false, "This message.")
)

type Action struct {
	Name    string
	Payload interface{}
}

func main() {
	flag.Parse()
	if *help || *processFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	process, err := process.ParseFile(*processFile)
	if err != nil {
		log.Fatalf("Error reading process definition: %s", err)
	}
	con := msg.Connect(*url)
	defer con.Close()
	sub := con.Subscribe(*topic)
	for {
		var action Action
		if sub.NextMessage(&action) {
			performAction(process, &action)
		} else {
			log.Print("No message received. Trying again.")
		}
	}
}

func performAction(Process *process.Process, Action *Action) {
	actionDesc := Process.Workflow[Action.Name]
	handlerURL := Process.Handlers[actionDesc.Handler].URL
	json, err := json.Marshal(Action.Payload)
	if err != nil {
		log.Printf("ERROR marshalling JSON for action with name [%s]: %s", Action.Name, err)
	}
	log.Printf("Performing action with name [%s], handler [%s], handler URL [%s] and payload [%s]",
		Action.Name, actionDesc.Handler, handlerURL, json)
	_, err = http.Post(handlerURL, "application/json", bytes.NewReader(json))
	if err != nil {
		log.Printf("ERROR performing action with name [%s], handler [%s], handler URL [%s]: %s",
			Action.Name, actionDesc.Handler, handlerURL, err)
	}
}
