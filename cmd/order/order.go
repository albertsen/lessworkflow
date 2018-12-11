package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/albertsen/lessworkflow/gen/proto/order"
)

var (
	url  = flag.String("u", "localhost:50051", "URL of Order Service.")
	help = flag.Bool("h", false, "This message.")
)

func main() {
	flag.Parse()
	command := flag.Arg(0)
	orderFile := flag.Arg(1)
	if *help || command != "place" || orderFile == "" {
		fmt.Printf("Usage: %s [options] place <order file>\n\nValid options:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	var order order.Order
	data, err := ioutil.ReadFile(orderFile)
	if err != nil {
		fmt.("Error reading file: %s", err)
	}
	json.Unmarshal(data, &order)

	// if err != nil {
	// 	log.Fatalf("Error reading file [%s]: %s", *messageFile, err)
	// }
	// con := msg.Connect(*url)
	// defer con.Close()
	// err = con.PublishBytes(*topic, message)
	// if err != nil {
	// 	log.Fatalf("Cannot publish message: %s", err)
	// }
}
