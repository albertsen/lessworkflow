package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/assert"
)

func ErrToGRPCStatusCode(t *testing.T, err error) codes.Code {
	stat, ok := status.FromError(err)
	assert.Assert(t, ok, "Error is not a gRPC error")
	return stat.Code()
}

func LoadTestData(Name string, TestData proto.Message) error {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return errors.New("GOPATH undefined")
	}
	file := path.Join(goPath, "/src/github.com/albertsen/lessworkflow/data/test/", Name) + ".json"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	json := string(data)
	return jsonpb.UnmarshalString(string(json), TestData)
}
