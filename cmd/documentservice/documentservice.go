package main

import (
	"context"

	ds "github.com/albertsen/lessworkflow/gen/proto/documentservice"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	dao "github.com/albertsen/lessworkflow/pkg/db/daos/documentdao"
	"github.com/albertsen/lessworkflow/pkg/grpc/grpcserver"
	"google.golang.org/grpc"
)

type service struct {
}

func (service *service) CreateDocument(ctx context.Context, req *ds.CreateDocumentRequest) (*ds.CreateDocumentResponse, error) {
	doc, err := dao.CreateDocment(req.Document)
	if err != nil {
		return &ds.CreateDocumentResponse{}, err
	}
	return &ds.CreateDocumentResponse{Document: doc}, nil
}

func (service *service) UpdateDocument(ctx context.Context, req *ds.UpdateDocumentRequest) (*ds.UpdateDocumentResponse, error) {
	doc, err := dao.UpdateDocument(req.Document)
	if err != nil {
		return &ds.UpdateDocumentResponse{}, err
	}
	return &ds.UpdateDocumentResponse{Document: doc}, nil
}

func (service *service) GetDocument(ctx context.Context, req *ds.GetDocumentRequest) (*ds.GetDocumentResponse, error) {
	doc, err := dao.GetDocument(req.Type, req.DocumentId)
	if err != nil {
		return &ds.GetDocumentResponse{}, err
	}
	return &ds.GetDocumentResponse{Document: doc}, nil
}

func (service *service) DeleteDocument(ctx context.Context, req *ds.DeleteDocumentRequest) (*ds.DeleteDocumentResponse, error) {
	err := dao.DeleteDocument(req.Type, req.DocumentId)
	if err != nil {
		return &ds.DeleteDocumentResponse{}, err
	}
	return &ds.DeleteDocumentResponse{}, nil

}

func main() {
	dbConn.Connect()
	defer dbConn.Close()
	grpcserver.StartServer(func(server *grpc.Server) {
		ds.RegisterDocumentServiceServer(server, &service{})
	})
}
