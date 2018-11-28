package event

import (
	"time"
)

type Action struct {
	Name       string
	RetryCount int
	DelayUntil time.Time
	Payload    interface{}
}
