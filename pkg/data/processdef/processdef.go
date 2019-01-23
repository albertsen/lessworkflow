package processdef

import (
	"fmt"
)

const (
	StepTypeAction = "action"
	StepTypeWait   = "wait"
)

type ProcessDef struct {
	Description string           `json:"description,omitempty"`
	Workflow    *ProcessWorkflow `json:"workflow"`
}

type ProcessWorkflow struct {
	Actions map[string]*ActionDef `json:"actions"`
	Steps   map[string]*StepDef   `json:"steps"`
	Start   string                `json:"start,omitempty"`
}

type ActionDef struct {
	URL string `json:"url"`
}

type StepDef struct {
	Action      string            `json:"action,omitempty"`
	Transitions map[string]string `json:"tansitions,omitempty"`
	WaitFor     string            `json:"waitFor,omitempty"`
	Next        string            `json:"next,omitempty"`
}

func (s *StepDef) Type() (string, error) {
	if s.Action == "" && s.WaitFor == "" {
		return "", fmt.Errorf("Invalid workflow step definition. Neither 'action' nor 'waitFor' attribute defined")
	}
	if s.Action != "" && s.WaitFor != "" {
		return "", fmt.Errorf("Invalid workflow step definition. Both 'action' and 'waitFor' attributes defined")
	}
	if s.Action != "" {
		return StepTypeAction, nil
	} else {
		return StepTypeWait, nil
	}
}

func (s *StepDef) ActionDef(pd *ProcessDef) (*ActionDef, error) {
	if pd == nil {
		return nil, fmt.Errorf("No process definition given")
	}
	if s.Action == "" {
		return nil, fmt.Errorf("Workflow step doesn't have an action defined")
	}
	if pd.Workflow == nil {
		return nil, fmt.Errorf("Process definition doesn't have a workflow defined")
	}
	if pd.Workflow.Actions == nil {
		return nil, fmt.Errorf("Process definition's workflow doesn't have any actions defind")
	}
	actionDef := pd.Workflow.Actions[s.Action]
	if actionDef == nil {
		return nil, fmt.Errorf("Cam't find action definotion for: %s", s.Action)
	}
	return actionDef, nil
}
