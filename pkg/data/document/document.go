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

func (doc *Document) GetContent(content interface{}) error {
	return json.Unmarshal(doc.Content, content)
}

func (doc *Document) SetContent(content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	err = doc.Content.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	return nil
}
