package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func main() {

	fmt.Println(os.Args)
	url := os.Args[1]
	subject := "hello"
	if len(os.Args) > 3 {
		subject = os.Args[3]
	}
	nc, err := nats.Connect(url)
	failOnError(err, "failed to connect to NATS")
	defer nc.Close()

	js, err := jetstream.New(nc)
	failOnError(err, "Failed to create a jetstream")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	failOnError(err, "Failed to define a ctx")
	defer cancel()

	streamName := fmt.Sprintf("%s_stream", subject)
	stream, err := js.Stream(ctx, streamName)
	failOnError(err, "Failed to create a stream")

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{})
	failOnError(err, "Failed to create a consumer")

	cc, err := cons.Consume(func(msg jetstream.Msg) {
		log.Printf("Received message: %s", string(msg.Data()))
		time.Sleep(1 * time.Second)
		msg.Ack()
	})
	failOnError(err, "Failed to create a consumer")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-cc.Closed()
}
