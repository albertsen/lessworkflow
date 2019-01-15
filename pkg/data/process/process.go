package process

import "time"

type Process struct {
	tableName     struct{}   `sql:"processes"`
	ID            string     `json:"id"`
	Status        string     `json:"status"`
	TimeCreated   *time.Time `json:"timeCreated"`
	TimeUpdated   *time.Time `json:"timeUpdated"`
	DocumentURL   string     `json:"documentURL"`
	ProcessDefURL string     `json:"processDefURL"`
}
