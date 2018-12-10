package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pbAction "github.com/albertsen/lessworkflow/gen/proto/action"
	"github.com/albertsen/lessworkflow/pkg/msg"
	"github.com/albertsen/lessworkflow/pkg/processdef"
)

var (
	url         = flag.String("s", "nats://localhost:4222", "URL of messaging server.")
	topic       = flag.String("t", "actions", "Message topic.")
	processFile = flag.String("p", "", "Process descriptor file. Mandatory.")
	help        = flag.Bool("h", false, "This message.")
)

func main() {
	flag.Parse()
	if *help || *processFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	processDef, err := processdef.ParseFile(*processFile)
	if err != nil {
		log.Fatalf("Error reading process definition: %s", err)
	}
	con, err := msg.Connect(*url)
	if err != nil {
		log.Fatalf("Error connecting to messaging server [%s]: %s", url, err)
	}
	defer con.Close()
	sub := con.Subscribe(*topic)
	for {
		var actionRequest pbAction.Request
		if sub.NextMessage(&actionRequest) {
			err := performAction(con, processDef, &actionRequest)
			if err != nil {
				log.Printf("ERROR in process [%s] - performing action [%s]: %s",
					actionRequest.ProcessId, actionRequest.Name, err)
			}
		} else {
			log.Print("No message received. Trying again.")
		}
	}
}

func performAction(Connection *msg.Connection, ProcessDef *processdef.ProcessDef, ActionRequest *pbAction.Request) error {
	actionDesc := ProcessDef.Workflow[ActionRequest.Name]
	handlerURL := ProcessDef.Handlers[actionDesc.Handler].URL
	log.Printf("Performing action: process [%s] - action [%s] - handler [%s] - handler URL [%s]",
		ActionRequest.ProcessId, ActionRequest.Name, actionDesc.Handler, handlerURL)
	jsonDoc, err := json.Marshal(ActionRequest.Payload.Content)
	if err != nil {
		return err
	}
	resp, err := http.Post(handlerURL, "application/json", bytes.NewReader(jsonDoc))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var actionResponse pbAction.Response
	err = json.Unmarshal(body, &actionResponse)
	if err != nil {
		return err
	}
	log.Printf("Result of action: process [%s] - action [%s]: %s",
		ActionRequest.ProcessId, ActionRequest.Name, actionResponse.Result)
	if actionDesc.Transitions == nil {
		log.Printf("No further transition: process [%s] - action [%s]", ActionRequest.ProcessId, ActionRequest.Name)
		return nil
	}
	nextAction := actionDesc.Transitions[actionResponse.Result]
	if nextAction == "" {
		return fmt.Errorf("Cannot find transition for result: %s", actionResponse.Result)
	}
	var nextActionRequest = pbAction.Request{
		Name:       nextAction,
		RetryCount: 0,
		Payload:    actionResponse.Payload,
		ProcessId:  ActionRequest.ProcessId,
	}
	log.Printf("Requesting action: process [%s] - action: %s", nextActionRequest.ProcessId, nextActionRequest.Name)
	Connection.PublishProtobuf(*topic, &nextActionRequest)
	return nil
}
