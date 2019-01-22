package workflow

import (
	"time"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
)

type Step struct {
	Name       string        `json:"name"`
	ProcessID  string        `json:"processId"`
	RetryCount int32         `json:"retryCount"`
	DelayUtil  *time.Time    `json:"delayUntil"`
	ProcessDef *doc.Document `json:"processDef"`
	Document   *doc.Document `json:"document"`
}

type ActionRequest struct {
	Document *doc.Document `json:"document"`
}

type ActionResponse struct {
	Result   string        `json:"result"`
	Document *doc.Document `json:"document"`
}
