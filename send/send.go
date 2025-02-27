package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
	}
}

func main() {

	fmt.Print(os.Args)
	url := os.Args[1]
	messageCount, err := strconv.Atoi(os.Args[2])
	subject := "hello"
	if len(os.Args) > 3 {
		subject = os.Args[3]
	}
	failOnError(err, "Failed to parse second args as messageCount : int")
	nc, err := nats.Connect(url)
	failOnError(err, "Failed to connect to NATS")
	defer nc.Close()

	js, err := jetstream.New(nc)
	failOnError(err, "Failed to create a jetstream")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	failOnError(err, "Failed to define a ctx")
	defer cancel()

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     fmt.Sprintf("%s_stream", subject),
		Subjects: []string{subject},
	})
	failOnError(err, "Failed to create a stream")

	for i := 0; i < messageCount; i++ {
		body := fmt.Sprintf("Hello World: %d", i)
		_, err := js.Publish(
			ctx, subject, []byte(body),
		)
		log.Printf(" [x] Sent %s", body)
		failOnError(err, "Failed to publish a message")
	}
}
