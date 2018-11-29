package order

import (
	"fmt"
	"time"
)

type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type Order struct {
	ID          string     `json:"id"`
	TimePlaced  time.Time  `json:"timePlaced,string"`
	TimeCreated time.Time  `json:"timeCreated,string"`
	Version     int        `json:"version"`
	Status      string     `json:"status"`
	Total       Money      `json:"total"`
	LineItems   []LineItem `json:"lineItems"`
}

type LineItem struct {
	ProductID         string `json:"productID"`
	PrudctDescription string `json:"productDescription"`
	Count             int    `json:"count"`
	ItemPrice         Money  `json:"itemPrice"`
	Total             Money  `json:"total"`
}

func (m *Money) String() string {
	v := float32(m.Amount) / 10
	return fmt.Sprintf("%f %s", v, m.Currency)
}
