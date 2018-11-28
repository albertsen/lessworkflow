package process

type Handler struct {
	URL string
}

type Action struct {
	Handler     string
	Transitions map[string]string
}

type Process struct {
	Handlers map[string]Handler
	Workflow map[string]Action
}
