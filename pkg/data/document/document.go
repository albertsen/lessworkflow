package document

import (
	"encoding/json"
	"time"
)

type Document struct {
	tableName   struct{}        `sql:"documents"`
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	TimeCreated *time.Time      `json:"timeCreated"`
	TimeUpdated *time.Time      `json:"timeUpdated"`
	Version     int32           `json:"version"`
	Content     json.RawMessage `json:"content"`
}
