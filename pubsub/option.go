package pubsub

import "time"

const (
	// DefaultAckTimeout is the time we tell pubsub to wait for us to process a message before it requeues
	defaultAckTimeout = 60 * time.Second
	// Subscription retry policy, see: pubsub.RetryPolicy
	// Those consts may need to be tuned
	defaultMinimumBackoff      = 10 * time.Second
	defaultMaximumBackoff      = 600 * time.Second
	defaultExpirationPolicy    = 0 * time.Second
	defaultRetainAckedMessages = false
	defaultRetentionDuration   = 3600 * time.Second
)

type Option interface {
	apply(*options)
}

// config implements Config
type options struct {
	// TBD: Add more fields for customized settings, and set some more default values for them
	ackTimeout          time.Duration
	minimumBackoff      time.Duration
	maximumBackoff      time.Duration
	labels              map[string]string
	expirationPolicy    time.Duration
	retainAckedMessages bool
	retentionDuration   time.Duration
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) { f(o) }

// NewConfig creates a new pubsub config
func NewConfig(opts ...Option) *options {
	options := options{
		ackTimeout:          defaultAckTimeout,
		minimumBackoff:      defaultMinimumBackoff,
		maximumBackoff:      defaultMaximumBackoff,
		expirationPolicy:    defaultExpirationPolicy,
		retainAckedMessages: defaultRetainAckedMessages,
		retentionDuration:   defaultRetentionDuration,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *Config as the argument
		opt.apply(&options)
	}

	return &options
}

func WithAckTimeout(ackTimeout time.Duration) Option {
	return optionFunc(func(opt *options) {
		opt.ackTimeout = ackTimeout
	})
}

func WithMinimumBackoff(minimumBackoff time.Duration) Option {
	return optionFunc(func(opt *options) {
		opt.minimumBackoff = minimumBackoff
	})
}

func WithMaximumBackoff(maximumBackoff time.Duration) Option {
	return optionFunc(func(opt *options) {
		opt.maximumBackoff = maximumBackoff
	})
}

func WithLabels(labels map[string]string) Option {
	return optionFunc(func(opt *options) {
		opt.labels = labels
	})
}

func WithExpirationPolicy(expirationPolicy time.Duration) Option {
	return optionFunc(func(opt *options) {
		opt.expirationPolicy = expirationPolicy
	})
}

func WithRetainAckedMessages() Option {
	return optionFunc(func(opt *options) {
		opt.retainAckedMessages = true
	})
}

func WithoutRetainAckedMessages() Option {
	return optionFunc(func(opt *options) {
		opt.retainAckedMessages = false
	})
}

func WithRetentionDuration(retentionDuration time.Duration) Option {
	return optionFunc(func(opt *options) {
		opt.retentionDuration = retentionDuration
	})
}
