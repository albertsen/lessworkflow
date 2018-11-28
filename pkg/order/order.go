package order

import (
	"fmt"
	"time"
)

type Money struct {
	Amount   int64
	Currency string
}

type Order struct {
	ID          string
	TimePlaced  time.Time
	TimeUpdated time.Time
	Version     int
	Status      string
	Total       Money
	LineItems   []LineItem
}

type LineItem struct {
	ProductID         string
	PrudctDescription string
	Count             int
	ItemPrice         Money
	Total             Money
}

func (m *Money) String() string {
	v := float32(m.Amount) / 10
	return fmt.Sprintf("%f %s", v, m.Currency)
}
