package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/albertsen/lessworkflow/pkg/rest/server"
	"github.com/albertsen/lessworkflow/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var (
	dbClient *mongo.Client
	dbName   = utils.Getenv("DB_NAME", "lessworkflow")
	router   = mux.NewRouter()
)

type RestHandler struct {
	Handler         func(doc interface{}, params map[string]string) *RestResponse
	DocumentFactory func() interface{}
}

func (rh *RestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var doc interface{}
	if rh.DocumentFactory != nil {
		doc = rh.DocumentFactory()
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			server.SendError(w, http.StatusUnprocessableEntity, err)
			return
		}
	}
	res := rh.Handler(doc, mux.Vars(r))
	server.SendResponse(w, res.StatusCode, res.Body)

}

type RestResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       interface{}
}

type RestError struct {
	Message string `json:"message"`
}

func NewErrorRestResponse(statusCode int, message interface{}) *RestResponse {
	return &RestResponse{
		StatusCode: statusCode,
		Body: &RestError{
			Message: fmt.Sprintf("%s", message),
		},
	}
}

func MapFactory() interface{} {
	return make(map[string]interface{})
}

func Handle(path string,
	handleFunc func(interface{}, map[string]string) *RestResponse,
	docFactory func() interface{}) *mux.Route {
	return router.Handle(path, &RestHandler{
		Handler:         handleFunc,
		DocumentFactory: docFactory,
	})
}

func DBCollection(collection string) *mongo.Collection {
	return dbClient.Database(dbName).Collection(collection)
}

func CreateDocument(v interface{}, params map[string]string) *RestResponse {
	if v == nil {
		return NewErrorRestResponse(http.StatusUnprocessableEntity, "Cannot create empty document")
	}
	doc, ok := v.(map[string]interface{})
	if !ok {
		return NewErrorRestResponse(http.StatusUnprocessableEntity, "Cannot convert document to map")
	}
	now := time.Now()
	doc["timeCreated"] = &now
	doc["timeUpdated"] = &now
	resource := params["resource"]
	res, err := DBCollection(resource).InsertOne(context.Background(), doc)
	if err != nil {
		return NewErrorRestResponse(http.StatusInternalServerError, err)
	}
	doc["_id"] = res.InsertedID
	return &RestResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Location": fmt.Sprintf("/%s/%s", resource, res.InsertedID)},
		Body:       doc,
	}
}

func UpdateDocument(v interface{}, params map[string]string) *RestResponse {
	if v == nil {
		return NewErrorRestResponse(http.StatusUnprocessableEntity, "Cannot update with empty document")
	}
	doc, ok := v.(map[string]interface{})
	if !ok {
		return NewErrorRestResponse(http.StatusUnprocessableEntity, "Cannot convert document to map")
	}
	now := time.Now()
	doc["timeUpdated"] = &now
	resource := params["resource"]
	id := params["id"]
	if id != doc["_id"] {
		return NewErrorRestResponse(http.StatusNotFound, fmt.Sprintf("ID in URL [%s] does not match ID in document [%s]", id, doc["_id"]))
	}
	res := DBCollection(resource).FindOneAndReplace(context.Background(), bson.M{"_id": id}, doc)
	if res.Err() != nil {
		return NewErrorRestResponse(http.StatusInternalServerError, res.Err())
	}
	var updatedDoc interface{}
	err := res.Decode(&updatedDoc)
	if err != nil {
		return NewErrorRestResponse(http.StatusInternalServerError, res.Err())
	}
	return &RestResponse{
		StatusCode: http.StatusOK,
		Body:       updatedDoc,
	}
}

func GetDocument(doc interface{}, params map[string]string) *RestResponse {
	resource := params["resource"]
	id := params["id"]
	res := DBCollection(resource).FindOne(context.Background(), bson.M{"_id": id})
	if res.Err() != nil {
		return NewErrorRestResponse(http.StatusInternalServerError, res.Err())
	}
	var dbDoc interface{}
	err := res.Decode(&dbDoc)
	if err != nil {
		return NewErrorRestResponse(http.StatusInternalServerError, err)
	}
	return &RestResponse{
		StatusCode: http.StatusOK,
		Body:       dbDoc,
	}
}

func DeleteDocument(doc interface{}, params map[string]string) *RestResponse {
	resource := params["resource"]
	id := params["id"]
	res, err := DBCollection(resource).DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return NewErrorRestResponse(http.StatusInternalServerError, err)
	}
	if res.DeletedCount < 1 {
		return NewErrorRestResponse(http.StatusNotFound, fmt.Sprintf("Document not found /%s/%s", resource, id))
	}
	return &RestResponse{
		StatusCode: http.StatusOK,
	}
}

func main() {
	mc, err := mongo.NewClient(utils.Getenv("DB_URL", "mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = mc.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer mc.Disconnect(context.Background())
	dbClient = mc
	router.Handle(
		"/{resource}",
		&RestHandler{Handler: CreateDocument, DocumentFactory: MapFactory},
	).Methods("POST")
	router.Handle(
		"/{resource}/{id}",
		&RestHandler{Handler: UpdateDocument, DocumentFactory: MapFactory},
	).Methods("PUT")
	router.Handle(
		"/{resource}/{id}",
		&RestHandler{Handler: GetDocument},
	).Methods("GET")
	router.Handle(
		"/{resource}/{id}",
		&RestHandler{Handler: CreateDocument},
	).Methods("DELETE")
	listenAddr := utils.Getenv("LISTEN_ADDR", ":8080")
	log.Printf("documentservice listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, router))
}
