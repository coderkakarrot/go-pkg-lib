package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pantheon-systems/go-pkg-lib/pubsub"
)

const (
	PROJECTID = "parallel-dynamic-runtime-tf"
	TOPICNAME = "test-pubsub-topic"
)

var log slog.Logger

func main() {
	ctx := context.Background()
	pubsubConfig := pubsub.NewConfig(
		pubsub.WithAckTimeout(300*time.Second),          // nolint:gomnd
		pubsub.WithMinimumBackoff(100*time.Millisecond), // nolint:gomnd
		pubsub.WithMaximumBackoff(60*time.Second),       // nolint:gomnd
	)
	c, err := pubsub.NewClient(ctx, PROJECTID, pubsubConfig)
	if err != nil {
		log.Error("could not create pubsub client", err)
		return
	}

	tp, err := c.Topic(ctx, TOPICNAME)
	if err != nil {
		log.Error("could not get pubsub topic", err)
		return
	}
	r, err := tp.Publish(ctx, &pubsub.Message{
		Data: []byte("Hello, World!"),
	})
	fmt.Println("Message ID: ", r)
}
