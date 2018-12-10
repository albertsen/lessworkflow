package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
)

const (
	topicName = "actions"
)

func main() {
	flag.Parse()
	command := flag.Arg(0)
	ctx := context.Background()
	projectID := "sap-se-commerce-arch"
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	topic := client.Topic(topicName)
	if command == "publish" {
		filename := flag.Arg(1)
		if filename == "" {
			printUsageAndExit()
		}
		err = publishFile(topic, filename)
		if err != nil {
			log.Fatalf("Failed to publish file: %s", err)
		}
	} else if command == "subscribe" {
		err = subscribe(client, topic)
		if err != nil {
			log.Fatalf("Failed to subscribe to topic: %s", err)
		}
	} else {
		printUsageAndExit()
	}
}

func printUsageAndExit() {
	fmt.Printf("Usage: %s publish <file>|subscribe\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func publishFile(Topic *pubsub.Topic, FileName string) error {
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		return err
	}
	ctx := context.Background()
	result := Topic.Publish(ctx, &pubsub.Message{
		Data: []byte(data),
	})
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}

func subscribe(Client *pubsub.Client, Topic *pubsub.Topic) error {
	ctx := context.Background()
	subName := "actions-subscription"
	var sub *pubsub.Subscription
	var err error
	sub = Client.Subscription(subName)
	if sub == nil {
		sub, err = Client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
			Topic:       Topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			return err
		}
	}
	cctx, cancel := context.WithCancel(ctx)
	fmt.Println("Receiving...")
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Got message: %q\n", string(msg.Data))
	})
	fmt.Println("Making channels...")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
	cancel()
	return err
}
