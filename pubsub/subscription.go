package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// subscription implements Pub/Sub subscription
type Subscription struct {
	*pubsub.Subscription
}

func (s *Subscription) Delete(ctx context.Context) error {
	// Check the subscription exists or not.
	found, err := s.Subscription.Exists(ctx)
	// Return an error, if the subscription does not exist.
	if err != nil {
		return fmt.Errorf("failed to check if subscription exists: %w", err)
	}

	if !found {
		return fmt.Errorf("subscription '%s' does not exist", s.ID())
	}

	if err := s.Subscription.Delete(ctx); err != nil {
		return fmt.Errorf("could not check if subscription '%s' exists: %w", s.ID(), err)
	}
	return nil
}
