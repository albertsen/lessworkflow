package main

import (
	"context"

	ops "github.com/albertsen/lessworkflow/gen/proto/orderprocessservice"
	"github.com/albertsen/lessworkflow/pkg/grpc/grpcserver"
	"google.golang.org/grpc"
)

type service struct {
}

func (service *service) StartProcess(ctx context.Context, req *ops.StartProcessRequest) (*ops.StartProcessResponse, error) {
	return nil, nil
}

func main() {
	grpcserver.StartServer(func(server *grpc.Server) {
		ops.RegisterOrderProcessServiceServer(server, &service{})
	})
}
