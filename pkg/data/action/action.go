package action

import (
	"time"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
)

type Action struct {
	Name       string        `json:"name"`
	ProcessID  string        `json:"processId"`
	RetryCount int32         `json:"retryCount"`
	DelayUtil  *time.Time    `json:"delayUntil"`
	ProcessDef *doc.Document `json:"processDef"`
	Document   *doc.Document `json:"document"`
}
