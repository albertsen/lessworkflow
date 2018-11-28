package process

import (
	"encoding/json"
)

func Parse(JSON []byte) (*Process, error) {
	var process Process
	err := json.Unmarshal(JSON, &process)
	return &process, err
}
