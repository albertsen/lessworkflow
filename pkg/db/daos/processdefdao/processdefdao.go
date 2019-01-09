package orderdao

import (
	"log"

	pd "github.com/albertsen/lessworkflow/gen/proto/processdef"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	"github.com/go-pg/pg"
	"github.com/gogo/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type processDefDTO struct {
	tableName   struct{} `sql:"process_defs"`
	ID          string
	Description string
	Workflow    string
}

func newProcessDefDTO(processDef *pd.ProcessDef) (*processDefDTO, error) {
	marshaller := jsonpb.Marshaler{}
	json, err := marshaller.MarshalToString(processDef.Workflow)
	if err != nil {
		return nil, err
	}
	return &processDefDTO{
		ID:          processDef.Id,
		Description: processDef.Description,
		Workflow:    json,
	}, nil
}

func newProcessDefProto(processDefDTO *processDefDTO) (*pd.ProcessDef, error) {
	var workflow pd.ProcessWorkflow
	if err := jsonpb.UnmarshalString(processDefDTO.Workflow, &workflow); err != nil {
		return nil, err
	}
	return &pd.ProcessDef{
		Id:          processDefDTO.ID,
		Description: processDefDTO.Description,
		Workflow:    &workflow,
	}, nil
}

func CreateProcessDef(processDef *pd.ProcessDef) (*pd.ProcessDef, error) {
	if processDef == nil {
		return nil, status.New(codes.InvalidArgument, "No process definition provided").Err()
	}
	if processDef.Id == "" {
		return nil, status.New(codes.InvalidArgument, "No procesess definition ID provided").Err()
	}
	processDefDTO, err := newProcessDefDTO(processDef)
	if err = dbConn.DB().Insert(processDefDTO); err != nil {
		log.Printf("ERROR inserting process definition into DB: %s", err)
		return nil, err
	}
	return processDef, nil
}

func UpdateProcessDef(processDef *pd.ProcessDef) (*pd.ProcessDef, error) {
	if processDef.Id == "" {
		return nil, status.New(codes.InvalidArgument, "Process definition is missing ID").Err()
	}
	processDefDTO, err := newProcessDefDTO(processDef)
	if err != nil {
		log.Printf("ERROR converting protbuf process definition into ProcessDefDTO: %s", err)
		return nil, err
	}
	if err := dbConn.DB().Update(processDefDTO); err != nil {
		if err == pg.ErrNoRows {
			return nil, status.New(codes.NotFound, "Process definition not found").Err()
		}
		log.Printf("ERROR updating process definition in DB: %s", err)
		return nil, err
	}
	return processDef, nil
}

func GetProcessDef(processDefID string) (*pd.ProcessDef, error) {
	if processDefID == "" {
		return nil, status.New(codes.InvalidArgument, "No process definition ID provided").Err()
	}
	processDefDTO := &processDefDTO{
		ID: processDefID,
	}
	if err := dbConn.DB().Select(processDefDTO); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		log.Printf("ERROR selecting process definition from DB: %s", err)
		return nil, err
	}
	processDef, err := newProcessDefProto(processDefDTO)
	if err != nil {
		log.Printf("ERROR converting process definition DTO to process defintiom proto: %s", err)
		return nil, err
	}
	return processDef, nil
}

func DeleteProcessDef(processDefID string) error {
	if processDefID == "" {
		return status.New(codes.InvalidArgument, "No process definition ID provided").Err()
	}
	var processDefDTO = processDefDTO{
		ID: processDefID,
	}
	if err := dbConn.DB().Delete(&processDefDTO); err != nil {
		log.Printf("ERROR deleting process definition from DB: %s", err)
		return err
	}
	return nil
}
