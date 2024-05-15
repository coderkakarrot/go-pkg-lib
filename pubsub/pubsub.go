package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// PubSub is a wrapper around the pubsub client
type PubSub struct {
	Client    *pubsub.Client
	projectID string
	config    *options
}

// NewClient creates a new PubSub client
func NewClient(ctx context.Context, projectID string, config *options) (*PubSub, error) {
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not create pubsub client: %w", err)
	}

	return &PubSub{
		Client:    c,
		projectID: projectID,
		config:    config,
	}, nil
}

// Topic returns a reference to the existing pubsub topic if the topic ID exists.
func (p *PubSub) Topic(ctx context.Context, topicID string) (*Topic, error) {
	if topicID == "" {
		return nil, fmt.Errorf("invalid topicID: '%s'", topicID)
	}
	tp := p.Client.Topic(topicID)
	found, err := tp.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not check if topic '%s' exists: %w", topicID, err)
	}

	if found {
		return &Topic{tp}, nil
	}
	return nil, fmt.Errorf("topic '%s' does not exist", topicID)
}

// NewTopic returns a reference to the existing pubsub topic or creates a new one
// if the given topic ID does not exist.
func (p *PubSub) NewTopic(ctx context.Context, topicID string) (*Topic, error) {
	if topicID == "" {
		return nil, fmt.Errorf("invalid topicID: '%s'", topicID)
	}

	tp := p.Client.Topic(topicID)
	found, err := tp.Exists(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not check if topic '%s' exists: %w", topicID, err)
	}

	if found {
		// Return the existing topic if found
		return &Topic{
			tp,
		}, nil
	}
	// Create a new topic if not found
	tp, err = p.Client.CreateTopic(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("topic initialization error: %w", err)
	}

	return &Topic{
		tp,
	}, nil
}

// NewSubscription returns a reference to the existing pubsub subscription or creates a new one
// if the given subscription ID does not exist.
func (p *PubSub) NewSubscription(ctx context.Context, subID string, topicID string) (*Subscription, error) {
	cfg := p.config

	if subID == "" {
		return nil, fmt.Errorf("invalid subscriptionID: '%s'", subID)
	}

	sub := p.Client.Subscription(subID)
	found, err := sub.Exists(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not check if subscription '%s' exists: %w", subID, err)
	}

	// Return the existing topic if found
	if found {
		return &Subscription{
			sub,
		}, nil
	}

	// Create a new topic if not found
	sub, err = p.Client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic:               p.Client.Topic(topicID),
		AckDeadline:         cfg.ackTimeout,
		RetainAckedMessages: cfg.retainAckedMessages,
		RetentionDuration:   cfg.retentionDuration,
		ExpirationPolicy:    cfg.expirationPolicy,
		Labels:              cfg.labels,
		RetryPolicy: &pubsub.RetryPolicy{
			MinimumBackoff: cfg.minimumBackoff,
			MaximumBackoff: cfg.maximumBackoff,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("subscription initialization error: %w", err)
	}

	return &Subscription{
		sub,
	}, nil
}

// Subscription returns a reference to the existing pubsub subscription if the subscription ID exists.
func (p *PubSub) Subscription(ctx context.Context, subID string) (*Subscription, error) {
	sub := p.Client.Subscription(subID)
	if sub == nil {
		return nil, fmt.Errorf("subscription '%s' does not exist", subID)
	}

	found, err := sub.Exists(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not check if subscription '%s' exists: %w", subID, err)
	}

	if found {
		return &Subscription{
			sub,
		}, nil
	}

	return nil, fmt.Errorf("subscription '%s' does not exist", subID)
}

// ProjectID is an accessor for the project ID
func (p *PubSub) ProjectID() string { return p.projectID }
