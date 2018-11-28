package main

import (
	"flag"
	"log"

	"github.com/albertsen/lw-processengine/cmd/common"
	"github.com/albertsen/lw-processengine/pkg/msg"
	"github.com/albertsen/lw-processengine/pkg/process"
)

var (
	url, topic, help      = common.Flags()
	processDefinitionFile = flag.String("p", "", "Process definition file. Will be read from stdin if ommitted.")
)

type Action struct {
	Name    string
	Payload interface{}
}

func main() {
	common.ParseFlags(help)
	json, err := common.ReadFileOrStdin(*processDefinitionFile)
	if err != nil {
		log.Fatalf("Error reading process definition file: %s", err)
	}
	process, err := process.Parse(json)
	if err != nil {
		log.Fatalf("Error parsing process definition: %s", err)
	}
	con := msg.Connect(*url)
	defer con.Close()
	sub := con.Subscribe(*topic)
	for {
		var action Action
		if sub.NextMessage(&action) {
			performAction(process, &action)
		} else {
			log.Print("No message received. Trying again.")
		}
	}
}

func performAction(Process *process.Process, Action *Action) {
	log.Printf("Action name: %s", Action.Name)
	log.Printf("Action payload: %s", Action.Payload)
	actionDesc := Process.Workflow[Action.Name]
	log.Printf("Action handler: %s", actionDesc.Handler)
	handlerURL := Process.Handlers[actionDesc.Handler].URL
	log.Printf("Action handler URL: %s", handlerURL)
}
