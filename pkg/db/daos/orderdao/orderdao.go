package orderdao

import (
	"log"
	"time"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	dbConn "github.com/albertsen/lessworkflow/pkg/db/conn"
	putils "github.com/albertsen/lessworkflow/pkg/proto/utils"
	"github.com/go-pg/pg"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderDTO struct {
	tableName   struct{} `sql:"orders"`
	ID          string
	TimeCreated *time.Time
	TimeUpdated *time.Time
	TimePlaced  *time.Time
	Version     int32
	Status      string
	Details     string
}

func newOrderDTO(order *od.Order) (*orderDTO, error) {
	marshaller := jsonpb.Marshaler{}
	json, err := marshaller.MarshalToString(order.Details)
	if err != nil {
		return nil, err
	}
	return &orderDTO{
		ID:          order.Id,
		TimeCreated: putils.NewTime(order.TimeCreated),
		TimeUpdated: putils.NewTime(order.TimeUpdated),
		TimePlaced:  putils.NewTime(order.TimeCreated),
		Version:     order.Version,
		Status:      order.Status,
		Details:     json,
	}, nil
}

func newOrderProto(orderDTO *orderDTO) (*od.Order, error) {
	var orderDetails od.OrderDetails
	if err := jsonpb.UnmarshalString(orderDTO.Details, &orderDetails); err != nil {
		return nil, err
	}
	return &od.Order{
		Id:          orderDTO.ID,
		TimeCreated: putils.NewTimestamp(orderDTO.TimeCreated),
		TimeUpdated: putils.NewTimestamp(orderDTO.TimeUpdated),
		TimePlaced:  putils.NewTimestamp(orderDTO.TimePlaced),
		Version:     orderDTO.Version,
		Status:      orderDTO.Status,
		Details:     &orderDetails,
	}, nil
}

func CreateOrder(order *od.Order) (*od.Order, error) {
	if order == nil {
		return nil, status.New(codes.InvalidArgument, "No order provided").Err()
	}
	if order.Id != "" {
		return nil, status.New(codes.InvalidArgument, "New order cannot have an ID").Err()
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Printf("ERROR generating UUID: %s", err)
		return nil, err
	}
	order.Id = uuid.String()
	now := ptypes.TimestampNow()
	order.TimeCreated = now
	order.TimeUpdated = now
	if order.TimePlaced == nil {
		order.TimePlaced = now
	}
	order.Status = "CREATED"
	order.Version = 1
	orderDTO, err := newOrderDTO(order)
	if err != nil {
		log.Printf("ERROR converting protbuf Order into OrderDTO: %s", err)
		return nil, err
	}
	if err = dbConn.DB().Insert(orderDTO); err != nil {
		log.Printf("ERROR inserting Order into DB: %s", err)
		return nil, err
	}
	return order, nil
}

func UpdateOrder(order *od.Order) (*od.Order, error) {
	if order.Id == "" {
		return nil, status.New(codes.InvalidArgument, "Order is missing ID").Err()
	}
	order.TimeUpdated = ptypes.TimestampNow()
	orderDTO, err := newOrderDTO(order)
	if err != nil {
		log.Printf("ERROR converting protbuf Order into OrderDTO: %s", err)
		return nil, err
	}
	if err := dbConn.DB().Update(orderDTO); err != nil {
		if err == pg.ErrNoRows {
			return nil, status.New(codes.NotFound, "Order not found").Err()
		}
		log.Printf("ERROR updating Order in DB: %s", err)
		return nil, err
	}
	return order, nil
}

func GetOrder(orderID string) (*od.Order, error) {
	if orderID == "" {
		return nil, status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
	orderDTO := &orderDTO{
		ID: orderID,
	}
	if err := dbConn.DB().Select(orderDTO); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		log.Printf("ERROR selecting Order from DB: %s", err)
		return nil, err
	}
	order, err := newOrderProto(orderDTO)
	if err != nil {
		log.Printf("ERROR converting Order DTO to Order proto: %s", err)
		return nil, err
	}
	return order, nil
}

func DeleteOrder(orderID string) error {
	if orderID == "" {
		return status.New(codes.InvalidArgument, "No order ID provided").Err()
	}
	var orderDTO = orderDTO{
		ID: orderID,
	}
	if err := dbConn.DB().Delete(&orderDTO); err != nil {
		log.Printf("ERROR deleting Order from DB: %s", err)
		return err
	}
	return nil
}
