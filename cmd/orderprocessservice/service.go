package main

import (
	"context"

	ops "github.com/albertsen/lessworkflow/gen/proto/orderprocessservice"
	"github.com/albertsen/lessworkflow/pkg/grpc/grpcserver"
	"google.golang.org/grpc"
)

type service struct {
}

func (service *service) DefineProcess(ctx context.Context, req *ops.DefineProcessRequest) (*ops.DefineProcessResponse, error) {
	return nil, nil
}

func (service *service) StartProcess(ctx context.Context, req *ops.StartProcessRequest) (*ops.StartProcessRequest, error) {
	return nil, nil
}

func main() {
	grpcserver.StartServer(func(server *grpc.Server) {
		ops.RegisterOrderProcessServiceServer(server, &service{})
	})
}
