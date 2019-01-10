package orderdao

import (
	"log"
	"time"

	doc "github.com/albertsen/lessworkflow/gen/proto/document"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	putils "github.com/albertsen/lessworkflow/pkg/proto/utils"
	"github.com/go-pg/pg"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type documentDTO struct {
	tableName   struct{} `sql:"documents"`
	ID          string
	Type        string
	TimeCreated *time.Time
	TimeUpdated *time.Time
	Version     int32
	Status      string
	Data        string
}

func newDocumentDTO(doc *doc.Document) (*documentDTO, error) {
	marshaller := jsonpb.Marshaler{}
	json, err := marshaller.MarshalToString(doc.Data)
	if err != nil {
		return nil, err
	}
	return &documentDTO{
		ID:          doc.Id,
		Type:        doc.Type,
		TimeCreated: putils.NewTime(doc.TimeCreated),
		TimeUpdated: putils.NewTime(doc.TimeUpdated),
		Version:     doc.Version,
		Data:        json,
	}, nil
}

func newDocumentProto(documentDTO *documentDTO) (*doc.Document, error) {
	var data any.Any
	if err := jsonpb.UnmarshalString(documentDTO.Data, &data); err != nil {
		return nil, err
	}
	return &doc.Document{
		Id:          documentDTO.ID,
		Type:        documentDTO.Type,
		TimeCreated: putils.NewTimestamp(documentDTO.TimeCreated),
		TimeUpdated: putils.NewTimestamp(documentDTO.TimeUpdated),
		Version:     documentDTO.Version,
		Data:        &data,
	}, nil
}

func CreateDocment(doc *doc.Document) (*doc.Document, error) {
	if doc == nil {
		return nil, status.New(codes.InvalidArgument, "Failed to create document: no document provided").Err()
	}
	if doc.Id == "" {
		return nil, status.New(codes.InvalidArgument, "Failed to create document: no ID provided").Err()
	}
	now := ptypes.TimestampNow()
	doc.TimeCreated = now
	doc.TimeUpdated = now
	doc.Version = 1
	docDTO, err := newDocumentDTO(doc)
	if err != nil {
		log.Printf("Failed to create document: Error converting protbuf into DTO - %s", err)
		return nil, err
	}
	if err = dbConn.DB().Insert(docDTO); err != nil {
		log.Printf("Failed to create document: Error inserting document into DB - %s", err)
		return nil, err
	}
	return doc, nil
}

func UpdateDocument(doc *doc.Document) (*doc.Document, error) {
	if doc == nil {
		return nil, status.New(codes.InvalidArgument, "Failed to update document: no document provided").Err()
	}
	if doc.Id == "" {
		return nil, status.New(codes.InvalidArgument, "Failed to update document: no ID provided").Err()
	}
	doc.TimeUpdated = ptypes.TimestampNow()
	docDTO, err := newDocumentDTO(doc)
	if err != nil {
		log.Printf("Failed to update document: Error converting protbuf into DTO - %s", err)
		return nil, err
	}
	if err := dbConn.DB().Update(docDTO); err != nil {
		if err == pg.ErrNoRows {
			return nil, status.New(codes.NotFound, "Order not found").Err()
		}
		log.Printf("Failed to update document: Error updating document in DB - %s", err)
		return nil, err
	}
	return doc, nil
}

func GetDocument(typeID string, docID string) (*doc.Document, error) {
	if typeID == "" {
		return nil, status.New(codes.InvalidArgument, "Failed to get document: No type provided").Err()
	}
	if docID == "" {
		return nil, status.New(codes.InvalidArgument, "Failed to get document: No document ID provided").Err()
	}
	documentDTO := documentDTO{
		ID:   docID,
		Type: typeID,
	}
	if err := dbConn.DB().Select(&documentDTO); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to get document: Error selecting document from DB - %s", err)
		return nil, err
	}
	order, err := newDocumentProto(&documentDTO)
	if err != nil {
		log.Printf("Failed to delete document: Error converting Order DTO to Order proto - %s", err)
		return nil, err
	}
	return order, nil
}

func DeleteDocument(typeID string, docID string) error {
	if typeID == "" {
		return status.New(codes.InvalidArgument, "Failed to delete document: No type provided").Err()
	}
	if docID == "" {
		return status.New(codes.InvalidArgument, "Failed to delete document: No document ID provided").Err()
	}
	documentDTO := documentDTO{
		ID:   docID,
		Type: typeID,
	}
	if err := dbConn.DB().Delete(&documentDTO); err != nil {
		log.Printf("Failed to delete document: Error deleting document from DB - %s", err)
		return err
	}
	return nil
}
