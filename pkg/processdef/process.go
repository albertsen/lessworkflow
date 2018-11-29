package processdef

type HandlerDef struct {
	URL string
}

type ActionDef struct {
	Handler     string
	Transitions map[string]string
}

type ProcessDef struct {
	Handlers map[string]HandlerDef
	Workflow map[string]ActionDef
}
