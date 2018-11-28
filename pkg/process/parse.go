package process

import (
	"encoding/json"
	"io/ioutil"
)

func ParseFile(File string) (*Process, error) {
	json, err := ioutil.ReadFile(File)
	if err != nil {
		return nil, err
	}
	return ParseJSON(json)
}

func ParseJSON(JSON []byte) (*Process, error) {
	var process Process
	err := json.Unmarshal(JSON, &process)
	return &process, err
}
