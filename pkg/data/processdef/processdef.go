package processdef

type ProcessDef struct {
	Description string           `json:"description,omitempty"`
	Workflow    *ProcessWorkflow `json:"workflow"`
}

type ProcessWorkflow struct {
	Handlers map[string]*HandlerDef `json:"handlers"`
	Actions  map[string]*ActionDef  `json:"actions"`
	Start    string                 `json:"start,omitempty"`
}

type HandlerDef struct {
	URL string `json:"url"`
}

type ActionDef struct {
	Handler     string            `json:"handler"`
	Transitions map[string]string `json:"tansitions"`
}
