package main

import (
	"context"

	pds "github.com/albertsen/lessworkflow/gen/proto/processdefservice"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	dao "github.com/albertsen/lessworkflow/pkg/db/daos/processdefdao"
	"github.com/albertsen/lessworkflow/pkg/grpc/grpcserver"
	"google.golang.org/grpc"
)

type service struct {
}

func (service *service) CreateProcessDef(ctx context.Context, req *pds.CreateProcessDefRequest) (*pds.CreateProcessDefResponse, error) {
	processDef, err := dao.CreateProcessDef(req.ProcessDef)
	if err != nil {
		return &pds.CreateProcessDefResponse{}, err
	}
	return &pds.CreateProcessDefResponse{ProcessDef: processDef}, nil
}

func (service *service) UpdateProcessDef(ctx context.Context, req *pds.UpdateProcessDefRequest) (*pds.UpdateProcessDefResponse, error) {
	processDef, err := dao.UpdateProcessDef(req.ProcessDef)
	if err != nil {
		return &pds.UpdateProcessDefResponse{}, err
	}
	return &pds.UpdateProcessDefResponse{ProcessDef: processDef}, nil
}

func (service *service) GetProcessDef(ctx context.Context, req *pds.GetProcessDefRequest) (*pds.GetProcessDefResponse, error) {
	processDef, err := dao.GetProcessDef(req.ProcessDefId)
	if err != nil {
		return &pds.GetProcessDefResponse{}, err
	}
	return &pds.GetProcessDefResponse{ProcessDef: processDef}, nil
}

func (service *service) DeleteProcessDef(ctx context.Context, req *pds.DeleteProcessDefRequest) (*pds.DeleteProcessDefResponse, error) {
	err := dao.DeleteProcessDef(req.ProcessDefId)
	if err != nil {
		return &pds.DeleteProcessDefResponse{}, err
	}
	return &pds.DeleteProcessDefResponse{}, nil

}

func main() {
	dbConn.Connect()
	defer dbConn.Close()
	grpcserver.StartServer(func(server *grpc.Server) {
		pds.RegisterProcessDefServiceServer(server, &service{})
	})
}
