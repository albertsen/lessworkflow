package main

import (
	"log"
	"time"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	ts "github.com/golang/protobuf/ptypes/timestamp"
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

func newTime(ts *ts.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	time, err := ptypes.Timestamp(ts)
	if err != nil {
		log.Printf("Error converting Time to protobuf Timestamp: %s", err)
		return nil
	}
	return &time
}

func newTimestamp(t *time.Time) *ts.Timestamp {
	if t == nil {
		return nil
	}
	ts, err := ptypes.TimestampProto(*t)
	if err != nil {
		log.Printf("Error converting protonf Timestamp to Time: %s", err)
		return nil
	}
	return ts
}

func newOrderDTO(order *od.Order) (*orderDTO, error) {
	marshaller := jsonpb.Marshaler{}
	json, err := marshaller.MarshalToString(order.Details)
	if err != nil {
		return nil, err
	}
	return &orderDTO{
		ID:          order.Id,
		TimeCreated: newTime(order.TimeCreated),
		TimeUpdated: newTime(order.TimeUpdated),
		TimePlaced:  newTime(order.TimeCreated),
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
		TimeCreated: newTimestamp(orderDTO.TimeCreated),
		TimeUpdated: newTimestamp(orderDTO.TimeUpdated),
		TimePlaced:  newTimestamp(orderDTO.TimePlaced),
		Version:     orderDTO.Version,
		Status:      orderDTO.Status,
		Details:     &orderDetails,
	}, nil
}
