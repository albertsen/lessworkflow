package utils

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	ts "github.com/golang/protobuf/ptypes/timestamp"
)

func NewTime(ts *ts.Timestamp) *time.Time {
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

func NewTimestamp(t *time.Time) *ts.Timestamp {
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
