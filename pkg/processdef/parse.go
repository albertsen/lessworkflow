package processdef

import (
	"encoding/json"
	"io/ioutil"
)

func ParseFile(File string) (*ProcessDef, error) {
	json, err := ioutil.ReadFile(File)
	if err != nil {
		return nil, err
	}
	return ParseJSON(json)
}

func ParseJSON(JSON []byte) (*ProcessDef, error) {
	var processDef ProcessDef
	err := json.Unmarshal(JSON, &processDef)
	return &processDef, err
}
