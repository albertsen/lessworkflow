package main

import (
	"fmt"
	"os"

	"github.com/albertsen/lessworkflow/pkg/msg"
)

func main() {
	err := msg.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// func performAction(Connection *msg.Connection, ProcessDef *processdef.ProcessDef, ActionRequest *pbAction.Request) error {
// 	actionDesc := ProcessDef.Workflow[ActionRequest.Name]
// 	handlerURL := ProcessDef.Handlers[actionDesc.Handler].URL
// 	log.Printf("Performing action: process [%s] - action [%s] - handler [%s] - handler URL [%s]",
// 		ActionRequest.ProcessId, ActionRequest.Name, actionDesc.Handler, handlerURL)
// 	jsonDoc, err := json.Marshal(ActionRequest.Payload.Content)
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := http.Post(handlerURL, "application/json", bytes.NewReader(jsonDoc))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	var actionResponse pbAction.Response
// 	err = json.Unmarshal(body, &actionResponse)
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("Result of action: process [%s] - action [%s]: %s",
// 		ActionRequest.ProcessId, ActionRequest.Name, actionResponse.Result)
// 	if actionDesc.Transitions == nil {
// 		log.Printf("No further transition: process [%s] - action [%s]", ActionRequest.ProcessId, ActionRequest.Name)
// 		return nil
// 	}
// 	nextAction := actionDesc.Transitions[actionResponse.Result]
// 	if nextAction == "" {
// 		return fmt.Errorf("Cannot find transition for result: %s", actionResponse.Result)
// 	}
// 	var nextActionRequest = pbAction.Request{
// 		Name:       nextAction,
// 		RetryCount: 0,
// 		Payload:    actionResponse.Payload,
// 		ProcessId:  ActionRequest.ProcessId,
// 	}
// 	log.Printf("Requesting action: process [%s] - action: %s", nextActionRequest.ProcessId, nextActionRequest.Name)
// 	Connection.PublishProtobuf(*topic, &nextActionRequest)
// 	return nil
// }
