package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

func LoadData(Path string, TestData interface{}) error {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return errors.New("GOPATH undefined")
	}
	file := path.Join(goPath, "/src/github.com/albertsen/lessworkflow/data/", Path)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, TestData)
}
