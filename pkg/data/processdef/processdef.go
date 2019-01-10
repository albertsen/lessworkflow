package processdef

type ProcessDef struct {
	ID          string
	Description string
	Workflow    *ProcessWorkflow
}

type ProcessWorkflow struct {
	Handlers map[string]*HandlerDef
	Actions  map[string]*ActionDef
}

type HandlerDef struct {
	UURL string
}

type ActionDef struct {
	Handler     string
	Transitions map[string]string
}
