package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// Topic is a wrapper around the pubsub topic
type Topic struct {
	*pubsub.Topic
}

// Message is a wrapper around the pubsub message
type Message = pubsub.Message

// Publish sends a message to the topic
// using the context provided
func (t *Topic) Publish(ctx context.Context, m *Message) (string, error) {
	if m == nil {
		return "", fmt.Errorf("message is empty")
	}

	r := t.Topic.Publish(ctx, m)
	if r == nil {
		return "", fmt.Errorf("failed to publish message")

	}

	return r.Get(ctx)
}

// Delete removes the topic from pubsub
// using the context provided
func (t *Topic) Delete(ctx context.Context) error {
	// Check if the topic exists or not.
	found, err := t.Topic.Exists(ctx)
	// Return an error, if the topic does not exist.
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %w", err)
	}
	if found {
		if err = t.Topic.Delete(ctx); err != nil {
			return fmt.Errorf("failed to delete topic: %w", err)
		}
	}
	return fmt.Errorf("topic '%s' does not exist", t.ID())
}
