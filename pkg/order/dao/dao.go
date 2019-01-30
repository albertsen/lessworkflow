package dao

import (
	"time"

	"github.com/albertsen/lessworkflow/pkg/db/repo"
	od "github.com/albertsen/lessworkflow/pkg/order/data"
)

func CreateOrder(order *od.Order) (int, *od.Order) {
	now := time.Now().Truncate(time.Microsecond)
	order.TimeCreated = &now
	order.TimeUpdated = &now
	order.Version = 1
	statusCode, err := repo.Insert(&order)
}
