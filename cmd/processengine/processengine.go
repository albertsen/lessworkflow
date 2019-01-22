package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pd "github.com/albertsen/lessworkflow/pkg/data/processdef"
	wf "github.com/albertsen/lessworkflow/pkg/data/workflow"
	"github.com/albertsen/lessworkflow/pkg/msg"
)

func newContentStruct() interface{} {
	return &wf.Step{}
}

func processContent(content interface{}) error {
	step, ok := content.(*wf.Step)
	if !ok {
		return errors.New(fmt.Sprintf("Unexepected content: %s", content))
	}
	return executeStep(step)
}

func main() {
	defer msg.Close()
	consumer := msg.NewConsumer("steps")
	done := make(chan bool)
	consumer.Consume(
		newContentStruct,
		processContent,
		done,
	)
	<-done
}

func executeStep(step *wf.Step) error {
	var processDef *pd.ProcessDef
	if step.Name == "" {
		return fmt.Errorf("Step without a name cannot be executed")
	}
	if step.ProcessDef == nil {
		return fmt.Errorf("No process definition attached to step: %s", step.Name)
	}
	if step.ProcessDef.Content == nil {
		return fmt.Errorf("Process defintion doesn't have content for step: %s", step.Name)
	}
	if err := json.Unmarshal(step.ProcessDef.Content, &processDef); err != nil {
		return err
	}
	if processDef.Workflow == nil {
		return fmt.Errorf("Process definition doesn't have workflow")
	}
	if processDef.Workflow.Steps == nil {
		return fmt.Errorf("Workflow doesn't have any steps")
	}
	stepDef := processDef.Workflow.Steps[step.Name]
	if stepDef == nil {
		return fmt.Errorf("Cannot find definition for step: %s", step.Name)
	}
	if stepDef.Action == "" {
		return fmt.Errorf("No action defined for step: %s", step.Name)
	}
	if processDef.Workflow.Actions == nil {
		return fmt.Errorf("Workflow doesn't have any actions")
	}

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
