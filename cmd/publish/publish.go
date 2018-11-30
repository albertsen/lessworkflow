package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/albertsen/lessworkflow/pkg/msg"
)

var (
	url         = flag.String("s", "nats://localhost:4222", "URL of messaging server.")
	topic       = flag.String("t", "actions", "Message topic.")
	messageFile = flag.String("m", "", "File containing the message to be published. Mandatory.")
	help        = flag.Bool("h", false, "This message.")
)

func main() {
	flag.Parse()
	if *help || *messageFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	message, err := ioutil.ReadFile(*messageFile)
	if err != nil {
		log.Fatalf("Error reading file [%s]: %s", *messageFile, err)
	}
	con := msg.Connect(*url)
	defer con.Close()
	err = con.PublishBytes(*topic, message)
	if err != nil {
		log.Fatalf("Cannot publish message: %s", err)
	}
}
