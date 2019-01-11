package action

import (
	"time"

	"github.com/albertsen/lessworkflow/pkg/data/processdef"
)

type Action struct {
	Name       string                 `json:"name"`
	ProcessID  string                 `json:"processId"`
	RetryCount int32                  `json:"retryCount"`
	DelayUtil  *time.Time             `json:"delayUntil"`
	ProcessDef *processdef.ProcessDef `json:"processDef"`
	Document   interface{}            `json:"document"`
}
