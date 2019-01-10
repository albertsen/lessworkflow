package orderdao

import (
	"log"
	"net/http"
	"time"

	doc "github.com/albertsen/lessworkflow/pkg/data/document"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateDocment(doc *doc.Document) (*doc.Document, error) {
	if doc == nil {
		return nil, status.New(codes.InvalidArgument, "Failed to create document: no document provided").Err()
	}
	if doc.ID != "" {
		return nil, status.New(codes.InvalidArgument, "Failed to create document: new document cannot have ID").Err()
	}
	if uuid, err := uuid.NewV4(); err != nil {
		httpError(w, http.StatusInternalServerError, err)
		return
	} else {
		order.ID = uuid.String()
	}
	now := time.Now()
	doc.TimeCreated = &now
	doc.TimeUpdated = &now
	doc.Version = 1
	if err := dbConn.DB().Insert(doc); err != nil {
		log.Printf("Failed to create document: Error inserting document into DB - %s", err)
		return nil, err
	}
	return doc, nil
}

func UpdateDocument(doc *doc.Document) (*doc.Document, error) {
	if doc == nil {
		return nil, status.New(codes.InvalidArgument, "Failed to update document: no document provided").Err()
	}
	if doc.ID == "" {
		return nil, status.New(codes.InvalidArgument, "Failed to update document: no ID provided").Err()
	}
	now := time.Now()
	doc.TimeUpdated = &now
	if err := dbConn.DB().Update(doc); err != nil {
		if err == pg.ErrNoRows {
			return nil, status.New(codes.NotFound, "Document not found").Err()
		}
		log.Printf("Failed to update document: Error updating document in DB - %s", err)
		return nil, err
	}
	return doc, nil
}

func GetDocument(Type string, DocID string) (*doc.Document, error) {
	if Type == "" {
		return nil, status.New(codes.InvalidArgument, "Failed to get document: No type provided").Err()
	}
	if DocID == "" {
		return nil, status.New(codes.InvalidArgument, "Failed to get document: No document ID provided").Err()
	}
	doc := &doc.Document{
		ID:   DocID,
		Type: Type,
	}
	if err := dbConn.DB().Select(doc); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to get document: Error selecting document from DB - %s", err)
		return nil, err
	}
	return doc, nil
}

func DeleteDocument(Type string, DocID string) error {
	if Type == "" {
		return status.New(codes.InvalidArgument, "Failed to delete document: No type provided").Err()
	}
	if DocID == "" {
		return status.New(codes.InvalidArgument, "Failed to delete document: No document ID provided").Err()
	}
	doc := doc.Document{
		ID:   DocID,
		Type: Type,
	}
	if err := dbConn.DB().Delete(doc); err != nil {
		log.Printf("Failed to delete document: Error deleting document from DB - %s", err)
		return err
	}
	return nil
}
