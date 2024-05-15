//go:build unit
// +build unit

package pubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	fakeAckTimeout        = 1 * time.Second
	fakeMinimumBackoff    = 2 * time.Second
	fakeMaximumBackoff    = 3 * time.Second
	fakeExpirationPolicy  = 4 * time.Second
	fakeRetentionDuration = 5 * time.Second
)

func TestNewConfigWithRetainAckedMessages(t *testing.T) {
	conf := NewConfig(
		WithAckTimeout(fakeAckTimeout),
		WithMinimumBackoff(fakeMinimumBackoff),
		WithMaximumBackoff(fakeMaximumBackoff),
		WithExpirationPolicy(fakeExpirationPolicy),
		WithRetainAckedMessages(),
		WithRetentionDuration(fakeRetentionDuration),
	)

	assert.Equal(t, fakeAckTimeout, conf.ackTimeout)
	assert.Equal(t, fakeMinimumBackoff, conf.minimumBackoff)
	assert.Equal(t, fakeMaximumBackoff, conf.maximumBackoff)
	assert.Equal(t, fakeExpirationPolicy, conf.expirationPolicy)
	assert.Equal(t, true, conf.retainAckedMessages)
	assert.Equal(t, fakeRetentionDuration, conf.retentionDuration)
}

func TestNewConfigWithoutRetainAckedMessages(t *testing.T) {
	conf := NewConfig(
		WithAckTimeout(fakeAckTimeout),
		WithMinimumBackoff(fakeMinimumBackoff),
		WithMaximumBackoff(fakeMaximumBackoff),
		WithExpirationPolicy(fakeExpirationPolicy),
		WithoutRetainAckedMessages(),
		WithRetentionDuration(fakeRetentionDuration),
	)

	assert.Equal(t, fakeAckTimeout, conf.ackTimeout)
	assert.Equal(t, fakeMinimumBackoff, conf.minimumBackoff)
	assert.Equal(t, fakeMaximumBackoff, conf.maximumBackoff)
	assert.Equal(t, fakeExpirationPolicy, conf.expirationPolicy)
	assert.Equal(t, false, conf.retainAckedMessages)
	assert.Equal(t, fakeRetentionDuration, conf.retentionDuration)
}
