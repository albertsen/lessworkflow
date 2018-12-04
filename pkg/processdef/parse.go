package processdef

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ParseFile(File string) (*ProcessDef, error) {
	log.Printf("Loading process definition from file: %s", File)
	json, err := ioutil.ReadFile(File)
	if err != nil {
		return nil, err
	}
	log.Printf("Process definition loaded: %s", json)
	return ParseJSON(json)
}

func ParseJSON(JSON []byte) (*ProcessDef, error) {
	var processDef ProcessDef
	err := json.Unmarshal(JSON, &processDef)
	return &processDef, err
}
