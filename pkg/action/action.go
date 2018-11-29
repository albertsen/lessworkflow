package action

import (
	"time"
)

type Request struct {
	Name       string    `json:"name"`
	ProcessID  string    `json:"processId"`
	RetryCount int       `json:"retryCount"`
	DelayUntil time.Time `json:"delayUntil,string"`
	Payload    Payload   `json:"payload"`
}

type Response struct {
	Result  string  `json:"result"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
