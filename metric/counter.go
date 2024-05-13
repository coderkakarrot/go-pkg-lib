package metric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// Counter is a metric that represents a single numerical value that only ever goes up.
type Counter struct {
	// Counter is the actual counter created.
	counter metric.Float64Counter
}

// NewCounter creates a new counter with the given name and description.
func (m *Metric) NewCounter(name, description string) (*Counter, error) {
	c, err := m.meter.Float64Counter(name, metric.WithDescription(description))
	if err != nil {
		return nil, err
	}

	return &Counter{
		counter: c,
	}, nil
}

// Add increments the given counter by the given value.
func (c *Counter) Add(ctx context.Context, incr float64, attr Attributes) {
	c.counter.Add(ctx, incr, metric.WithAttributes(attr.toOtel()...))
}
