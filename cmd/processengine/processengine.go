package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pd "github.com/albertsen/lessworkflow/pkg/data/processdef"
	pe "github.com/albertsen/lessworkflow/pkg/data/processexec"
	"github.com/albertsen/lessworkflow/pkg/msg"
)

var (
	publisher *msg.Publisher
)

func newContentStruct() interface{} {
	return &pe.Step{}
}

func processContent(content interface{}) error {
	step, ok := content.(*pe.Step)
	if !ok {
		return fmt.Errorf("Unexepected content: %s", content)
	}
	return executeStep(step)
}

func main() {
	msg.Connect()
	defer msg.Close()
	publisher = msg.NewPublisher("steps")
	consumer := msg.NewConsumer("steps")
	done := make(chan bool)
	consumer.Consume(
		newContentStruct,
		processContent,
		done,
	)
	<-done
}

func executeStep(step *pe.Step) error {
	stepDef, err := step.StepDef()
	if err != nil {
		return err
	}
	stepType, err := stepDef.Type()
	if err != nil {
		return err
	}
	if stepType == pd.StepTypeAction {
		return executeActionStep(step, stepDef)
	} else {
		log.Printf("Wait steps not implemented yet")
		return nil
	}
}

func executeActionStep(step *pe.Step, stepDef *pd.StepDef) error {
	var proecessDef pd.ProcessDef
	if step.ProcessDef == nil {
		return fmt.Errorf("Step doesn't have process definition document")
	}
	if step.ProcessDef.Content == nil {
		return fmt.Errorf("Process definition document doesn't have content")
	}
	if err := json.Unmarshal(step.ProcessDef.Content, &proecessDef); err != nil {
		return err
	}
	actionDef, err := stepDef.ActionDef(&proecessDef)
	if err != nil {
		return err
	}
	log.Printf("Performing action: process [%s] - process ID [%s] - step [%s] - action [%s] - action URL [%s]",
		step.ProcessDef.ID, step.ProcessID, step.Name, stepDef.Action, actionDef.URL)
	actionReq := pe.ActionRequest{Document: step.Document}
	jsonDoc, err := json.Marshal(actionReq)
	if err != nil {
		return err
	}
	resp, err := http.Post(actionDef.URL, "application/json", bytes.NewReader(jsonDoc))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var actionResponse pe.ActionResponse
	err = json.Unmarshal(body, &actionResponse)
	if err != nil {
		return err
	}
	if stepDef.Transitions == nil {
		log.Printf("ERROR - No further transitons for process [%s] - process ID [%s] - step [%s] - action [%s] - action URL [%s]",
			step.ProcessDef.ID, step.ProcessID, step.Name, stepDef.Action, actionDef.URL)
	}
	nextStepName := stepDef.Transitions[actionResponse.Result]
	if nextStepName == "" {
		return fmt.Errorf("Cannot find transition for result [%s] in process [%s]", actionResponse.Result, step.ProcessDef.ID)
	}
	var nextStep = pe.Step{
		ProcessID:  step.ProcessID,
		ProcessDef: step.ProcessDef,
		Name:       nextStepName,
		Document:   actionResponse.Document,
	}
	publisher.Publish(nextStep)
	return nil
}
