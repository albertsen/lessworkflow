package main

import (
	"fmt"

	pe "github.com/albertsen/lessworkflow/pkg/data/processexec"
	"github.com/albertsen/lessworkflow/pkg/msg"
)

var (
	publisher *msg.Publisher
)

func newContentStruct() interface{} {
	return &pe.Step{}
}

func processContent(content interface{}) error {
	step, ok := content.(*pe.Step)
	if !ok {
		return fmt.Errorf("Unexepected content: %s", content)
	}
	return step.Execute()
}

func main() {
	msg.Connect()
	defer msg.Close()
	publisher = msg.NewPublisher("steps")
	consumer := msg.NewConsumer("steps")
	done := make(chan bool)
	consumer.Consume(
		newContentStruct,
		processContent,
		done,
	)
	<-done
}
