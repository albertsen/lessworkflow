package main

import (
	"fmt"
	"time"

	od "github.com/albertsen/lessworkflow/gen/proto/orderdata"
	"github.com/golang/protobuf/jsonpb"
)

func main() {
	now := time.Now()
	order := od.Order{
		TimePlaced: &now,
		TotalPrice: &od.MonetaryAmount{
			Currency: "EUR",
			Value:    3250,
		},
		LineItems: []*od.LineItem{
			{
				ProductId:          "augustiner",
				ProductDescription: "Augustiner Hell",
				Count:              10,
				ItemPrice: &od.MonetaryAmount{
					Currency: "EUR",
					Value:    200,
				},
				TotalPrice: &od.MonetaryAmount{
					Currency: "EUR",
					Value:    2000,
				},
			},
			{
				ProductId:          "giesinger",
				ProductDescription: "Giesinger Erhellung",
				Count:              5,
				ItemPrice: &od.MonetaryAmount{
					Currency: "EUR",
					Value:    250,
				},
				TotalPrice: &od.MonetaryAmount{
					Currency: "EUR",
					Value:    1250,
				},
			},
		},
	}
	marshaller := jsonpb.Marshaler{}
	json, err := marshaller.MarshalToString(&order)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(json)
	}
}
