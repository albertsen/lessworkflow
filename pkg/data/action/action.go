package action

import (
	"time"

	"github.com/albertsen/lessworkflow/pkg/data/processdef"
)

type Action struct {
	Name       string
	ProcessID  string
	RetryCount int32
	DelayUtil  *time.Time
	ProcessDef *processdef.ProcessDef
	Payload    interface{}
}
