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

func (doc *Document) GetContent(C interface{}) error {
	return json.Unmarshal(doc.Content, C)
}

func (doc *Document) SetContent(C interface{}) error {
	data, err := json.Marshal(C)
	if err != nil {
		return err
	}
	err = doc.Content.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	return nil
}
