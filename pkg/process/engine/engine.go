package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/albertsen/lessworkflow/pkg/msg"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	pd "github.com/albertsen/lessworkflow/pkg/data/processdef"
	restClient "github.com/albertsen/lessworkflow/pkg/rest/client"
)

const (
	stepTypeAction = "action"
	stepTypeWait   = "wait"
)

type ProcessEngine struct {
	StepPublisher  *msg.Publisher
	ErrorPublisher *msg.Publisher
	StepConsumer   *msg.Consumer
}

func StartProcess(processDef *pd.ProcessDef) error {
	return nil
}

type StepExecutionRequest struct {
	Name       string        `json:"name"`
	ProcessID  string        `json:"processId"`
	RetryCount int32         `json:"retryCount"`
	DelayUtil  *time.Time    `json:"delayUntil"`
	ProcessDef *doc.Document `json:"processDef"`
	Document   *doc.Document `json:"document"`
}

func (s *StepExecutionRequest) Execute() (*StepExecutionRequest, error) {
	stepDef, err := s.stepDef()
	if err != nil {
		return nil, err
	}
	stepType, err := stepDefType(stepDef)
	if err != nil {
		return nil, err
	}
	if stepType == pd.StepTypeAction {
		return s.executeActionStep(stepDef)
	} else {
		log.Printf("Wait steps not implemented yet")
		return nil, nil
	}
}

func (s *StepExecutionRequest) executeActionStep(stepDef *pd.StepDef) (*StepExecutionRequest, error) {
	if s.ProcessDef == nil {
		return nil, fmt.Errorf("Step doesn't have process definition document")
	}
	if s.ProcessDef.Content == nil {
		return nil, fmt.Errorf("Process definition document doesn't have content")
	}
	var processDef pd.ProcessDef
	if err := json.Unmarshal(s.ProcessDef.Content, &processDef); err != nil {
		return nil, err
	}
	actionDef, err := actionDef(stepDef, &processDef)
	if err != nil {
		return nil, err
	}
	log.Printf("Performing action: process [%s] - process ID [%s] - step [%s] - action [%s] - action URL [%s]",
		s.ProcessDef.ID, s.ProcessID, s.Name, stepDef.Action, actionDef.URL)
	actionReq := ActionRequest{Document: s.Document}
	var actionResponse ActionResponse
	_, err = restClient.Post(actionDef.URL, actionReq, &actionResponse)
	if err != nil {
		return nil, err
	}
	if stepDef.Transitions == nil {
		log.Printf("ERROR - No further transitons for process [%s] - process ID [%s] - step [%s] - action [%s] - action URL [%s]",
			s.ProcessDef.ID, s.ProcessID, s.Name, stepDef.Action, actionDef.URL)
	}
	nextStepName := stepDef.Transitions[actionResponse.Result]
	if nextStepName == "" {
		return nil, fmt.Errorf("Cannot find transition for result [%s] in process [%s]", actionResponse.Result, s.ProcessDef.ID)
	}
	return &StepExecutionRequest{
		ProcessID:  s.ProcessID,
		ProcessDef: s.ProcessDef,
		Name:       nextStepName,
		Document:   actionResponse.Document,
	}, nil
}

func (s *StepExecutionRequest) stepDef() (*pd.StepDef, error) {
	if s.Name == "" {
		return nil, fmt.Errorf("Step without a name cannot be executed")
	}
	if s.ProcessDef == nil {
		return nil, fmt.Errorf("No process definition attached to step: %s", s.Name)
	}
	if s.ProcessDef.Content == nil {
		return nil, fmt.Errorf("Process defintion doesn't have content for step: %s", s.Name)
	}
	var processDef pd.ProcessDef
	if err := json.Unmarshal(s.ProcessDef.Content, &processDef); err != nil {
		return nil, err
	}
	if processDef.Workflow == nil {
		return nil, fmt.Errorf("Process definition doesn't have workflow")
	}
	if processDef.Workflow.Steps == nil {
		return nil, fmt.Errorf("Workflow doesn't have any steps")
	}
	stepDef := processDef.Workflow.Steps[s.Name]
	if stepDef == nil {
		return nil, fmt.Errorf("No workflow step definition found for step: %", s.Name)
	}
	return stepDef, nil
}

func stepDefType(stepDef *pd.StepDef) (string, error) {
	if stepDef == nil {
		return "", fmt.Errorf("No step definotion given")
	}
	if stepDef.Action == "" && stepDef.WaitFor == "" {
		return "", fmt.Errorf("Invalid workflow step definition. Neither 'action' nor 'waitFor' attribute defined")
	}
	if stepDef.Action != "" && stepDef.WaitFor != "" {
		return "", fmt.Errorf("Invalid workflow step definition. Both 'action' and 'waitFor' attributes defined")
	}
	if stepDef.Action != "" {
		return stepTypeAction, nil
	} else {
		return stepTypeWait, nil
	}
}

func actionDef(stepDef *pd.StepDef, pd *pd.ProcessDef) (*pd.ActionDef, error) {
	if stepDef == nil {
		return nil, fmt.Errorf("No step definotion given")
	}
	if pd == nil {
		return nil, fmt.Errorf("No process definition given")
	}
	if stepDef.Action == "" {
		return nil, fmt.Errorf("Workflow step doesn't have an action defined")
	}
	if pd.Workflow == nil {
		return nil, fmt.Errorf("Process definition doesn't have a workflow defined")
	}
	if pd.Workflow.Actions == nil {
		return nil, fmt.Errorf("Process definition's workflow doesn't have any actions defind")
	}
	actionDef := pd.Workflow.Actions[stepDef.Action]
	if actionDef == nil {
		return nil, fmt.Errorf("Cam't find action definotion for: %s", stepDef.Action)
	}
	return actionDef, nil
}

type ActionRequest struct {
	Document *doc.Document `json:"document"`
}

type ActionResponse struct {
	Result   string        `json:"result"`
	Document *doc.Document `json:"document"`
}
