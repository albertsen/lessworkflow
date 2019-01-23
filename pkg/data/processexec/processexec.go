package workflow

import (
	"encoding/json"
	"fmt"
	"time"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	pd "github.com/albertsen/lessworkflow/pkg/data/processdef"
)

type Step struct {
	Name       string        `json:"name"`
	ProcessID  string        `json:"processId"`
	RetryCount int32         `json:"retryCount"`
	DelayUtil  *time.Time    `json:"delayUntil"`
	ProcessDef *doc.Document `json:"processDef"`
	Document   *doc.Document `json:"document"`
}

func (s *Step) StepDef() (*pd.StepDef, error) {
	if s.Name == "" {
		return nil, fmt.Errorf("Step without a name cannot be executed")
	}
	if s.ProcessDef == nil {
		return nil, fmt.Errorf("No process definition attached to step: %s", step.Name)
	}
	if s.ProcessDef.Content == nil {
		return nil, fmt.Errorf("Process defintion doesn't have content for step: %s", step.Name)
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

}

type ActionRequest struct {
	Document *doc.Document `json:"document"`
}

type ActionResponse struct {
	Result   string        `json:"result"`
	Document *doc.Document `json:"document"`
}
