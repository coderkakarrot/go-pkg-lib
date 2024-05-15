//go:build unit
// +build unit

package pubsub

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testServer *pstest.Server = nil
var testTopic *Topic = nil

const (
	testProjectID = "parallel-dynamic-runtime-tf"
	testTopicID   = "test-pubsub-topic"
)

// TestMain sets up the test server and creates a test topic
func TestMain(m *testing.M) {
	if ok := setup(); ok != true {
		os.Exit(1)
	}

	code := m.Run()

	shutdown()
	os.Exit(code)
}

// setup starts the test server and creates a test topic
func setup() bool {
	// use the gcloud emulator
	testServer = pstest.NewServer()

	cfg := NewConfig()
	ctx := context.Background()
	client, err := NewClient(ctx, testProjectID, cfg)
	if err != nil {
		return false
	}

	testTopic, err = client.NewTopic(ctx, testTopicID)

	return err == nil
}

// shutdown stops the test server
func shutdown() {
	testTopic.Stop()
	testServer.Close()
}

func processNotifications(t *testing.T, ctx context.Context, sub Subscription) {
	for {
		select {
		case <-ctx.Done():
			err := sub.Delete(context.Background())
			assert.NoError(t, err)

			return
		default:
			err := sub.Receive(ctx, func(rCtx context.Context, msg *pubsub.Message) {
				t.Logf("new message from sub: %s", msg.ID)
				msg.Ack()
			})
			assert.NoError(t, err)
		}
	}
}

func publishMessages(t *testing.T) {
	msg := &Message{
		Data: []byte("Hello, World!\n"),
	}
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		id, err := testTopic.Publish(ctx, msg)
		t.Logf("message published: %s", id)
		assert.NoError(t, err)
	}
}

func hostName() string {
	hs, err := os.Hostname()
	if err != nil {
		hs = uuid.New().String()
	}

	return hs
}

func TestNewSubscriber(t *testing.T) {
	subID := hostName()

	ctx := context.Background()
	cfg := NewConfig()

	client, err := NewClient(ctx, testProjectID, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	tp, err := client.NewTopic(ctx, testTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	sub, err := client.NewSubscription(ctx, subID, testTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	go publishMessages(t)
	go processNotifications(t, ctx, *sub)
	time.Sleep(1 * time.Second)
}

func TestTopic(t *testing.T) {
	fakeTopicID := "fake-topic"

	ctx := context.Background()
	cfg := NewConfig()

	client, err := NewClient(ctx, testProjectID, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	tp, err := client.NewTopic(ctx, testTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	// reference to an existing topic
	tp2, err := client.Topic(ctx, testTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, tp2)
	assert.Equal(t, tp, tp2)

	// reference a topic that doesn't exist
	tp3, err := client.Topic(ctx, fakeTopicID)
	assert.Error(t, err)
	assert.Nil(t, tp3)
	assert.NotEqual(t, tp, tp3)
}

func TestSubscription(t *testing.T) {
	subID := hostName()
	fakeSubID := "fake-subscription"

	ctx := context.Background()
	cfg := NewConfig()

	client, err := NewClient(ctx, testProjectID, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	tp, err := client.NewTopic(ctx, testTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	sub, err := client.NewSubscription(ctx, subID, testTopicID)
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	// reference to an existing subscription
	sub2, err := client.Subscription(ctx, subID)
	assert.NoError(t, err)
	assert.NotNil(t, sub2)
	assert.Equal(t, sub, sub2)

	// reference a subscription that doesn't exist
	sub3, err := client.Subscription(ctx, fakeSubID)
	assert.Error(t, err)
	assert.Nil(t, sub3)
	assert.NotEqual(t, sub, sub3)
}
