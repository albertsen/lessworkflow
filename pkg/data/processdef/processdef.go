package processdef

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

func (s *StepDef) Type() string, error {
	if s.Action == "" && s.WaitFor == "" {
		return "", fmt.Errorf("Invalid step definition. Neither 'action' nor '")
	}
}
