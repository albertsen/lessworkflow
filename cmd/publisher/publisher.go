package main

import (
	"flag"
	"log"

	"github.com/albertsen/lw-processengine/cmd/common"
	"github.com/albertsen/lw-processengine/pkg/msg"
)

var (
	url, topic, help = common.Flags()
	messageFile      = flag.String("m", "", "File containing message to be published")
)

func main() {
	common.ParseFlags(help)
	json, err := common.ReadFileOrStdin(*messageFile)
	if err != nil {
		log.Fatalf("Error reading file [%s]: %s", *messageFile, err)
	}
	con := msg.Connect(*url)
	defer con.Close()
	err = con.PublishBytes(*topic, json)
	if err != nil {
		log.Fatalf("Cannot publish message: %s", err)
	}
}
